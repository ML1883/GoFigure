package analyzer

import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	//"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

// DistributionType represents different probability distributions
type DistributionType string

const (
	NormalDist      DistributionType = "normal"
	GammaDist       DistributionType = "gamma"
	BetaDist        DistributionType = "beta"
	ExponentialDist DistributionType = "exponential"
	LogNormalDist   DistributionType = "lognormal"
	EmpiricalDist   DistributionType = "empirical"
)

// DistributionParameters stores parameters for various probability distributions
type DistributionParameters struct {
	Type          DistributionType
	Mean          float64
	StdDev        float64
	Shape         float64 // Alpha for Gamma/Beta, unused for Normal
	Rate          float64 // Beta for Gamma, unused for Normal
	Scale         float64 // Used for LogNormal
	EmpiricalCDF  []float64
	Bins          []float64
	GoodnessOfFit float64 // Higher is better fit
}

// TextDistributionWithFitting extends TextDistributionModel with distribution fitting capabilities
type TextDistributionWithFitting struct {
	TextDistributionModel
	// Distribution type and parameters for each character
	CharDistributionType [36]DistributionParameters
	// Raw frequency data collected for each character across samples
	CharFrequencyData [36][]float64
}

// TextDistributionModel represents the statistical distribution of characters across multiple texts
type TextDistributionModel struct {
	// Character distributions for letters (a-z) and numbers (0-9)
	CharDistribution [36]float64
	// Standard deviation for each character distribution
	CharStdDev [36]float64
	// Position distributions (mean position of each character relative to text length)
	PositionDistribution [36]float64
	// Position standard deviation
	PositionStdDev [36]float64
	// Total samples used to build the model
	SampleCount int
	// Threshold for anomaly detection
	AnomalyThreshold float64
}

// CreateDistributionModel builds a statistical model from multiple text samples
// Note: This function assumes character frequencies follow a normal distribution
// and calculates mean and standard deviation for each character
func CreateDistributionModel(textSamples []string, anomalyThreshold float64) (*TextDistributionModel, error) {
	if len(textSamples) == 0 {
		return nil, fmt.Errorf("no text samples provided")
	}

	// Initialize the model
	model := &TextDistributionModel{
		SampleCount:      len(textSamples),
		AnomalyThreshold: anomalyThreshold,
	}

	// Collect letter data from all samples
	var allLetterData []*LetterData
	for _, text := range textSamples {
		letterData := AnalyzeLettersFromText(text)
		allLetterData = append(allLetterData, letterData)
	}

	// Calculate mean distribution for each character
	for i := 0; i < 36; i++ {
		var sum float64
		var positions []float64

		// Gather frequencies and positions for this character across all samples
		for _, ld := range allLetterData {
			// Calculate relative frequency (normalized by text length)
			if ld.TotalCount > 0 {
				relFreq := float64(ld.LetterNumberArray[i]) / float64(ld.TotalCount)
				sum += relFreq

				// Calculate relative positions (normalized by text length)
				for _, pos := range ld.PositionArray[i] {
					relPos := float64(pos) / float64(ld.TotalCount)
					positions = append(positions, relPos)
				}
			}
		}

		// Set mean frequency
		model.CharDistribution[i] = sum / float64(len(textSamples))

		// Calculate mean position if we have positions
		if len(positions) > 0 {
			var posSum float64
			for _, pos := range positions {
				posSum += pos
			}
			model.PositionDistribution[i] = posSum / float64(len(positions))
		}

		// Calculate standard deviations
		var freqSumSq float64
		for _, ld := range allLetterData {
			if ld.TotalCount > 0 {
				relFreq := float64(ld.LetterNumberArray[i]) / float64(ld.TotalCount)
				diff := relFreq - model.CharDistribution[i]
				freqSumSq += diff * diff
			}
		}

		// Set frequency standard deviation
		if len(textSamples) > 1 {
			model.CharStdDev[i] = math.Sqrt(freqSumSq / float64(len(textSamples)-1))
		}

		// Calculate position standard deviation
		if len(positions) > 1 {
			var posSumSq float64
			for _, pos := range positions {
				diff := pos - model.PositionDistribution[i]
				posSumSq += diff * diff
			}
			model.PositionStdDev[i] = math.Sqrt(posSumSq / float64(len(positions)-1))
		}
	}

	return model, nil
}

// SaveModel saves the distribution model to a file
func (m *TextDistributionModel) SaveModel(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(m)
	if err != nil {
		return err
	}

	return nil
}

// LoadModel loads a distribution model from a file
func LoadModel(filename string) (*TextDistributionModel, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var model TextDistributionModel
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// AnomalyScore calculates how different a text is from the distribution model
// Returns a score and anomalous character information
// Note: This method assumes character frequencies follow a normal distribution
// and uses z-scores (standard deviations from the mean) to detect anomalies
func (m *TextDistributionModel) AnomalyScore(text string) (float64, map[string]float64) {
	// Analyze the text
	letterData := AnalyzeLettersFromText(text)

	// Calculate z-scores for each character
	zScores := make(map[string]float64)
	var totalScore float64
	var significantDeviations int

	for i := 0; i < 36; i++ {
		// Skip characters with no distribution data
		if m.CharStdDev[i] == 0 {
			continue
		}

		// Calculate relative frequency for this character
		var relFreq float64
		if letterData.TotalCount > 0 {
			relFreq = float64(letterData.LetterNumberArray[i]) / float64(letterData.TotalCount)
		}

		// Calculate z-score (how many standard deviations from the mean)
		zScore := math.Abs(relFreq-m.CharDistribution[i]) / m.CharStdDev[i]

		// Only count significant deviations
		if zScore > 1.0 {
			var charLabel string
			if i < 10 {
				// It's a number
				charLabel = fmt.Sprintf("%d", i)
			} else {
				// It's a letter (a=10, b=11, etc.)
				charLabel = string(rune('a' + (i - 10)))
			}

			zScores[charLabel] = zScore
			totalScore += zScore
			significantDeviations++
		}
	}

	// Calculate position deviation
	//var positionScore float64
	for i := 0; i < 36; i++ {
		// Skip characters with no position data
		if len(letterData.PositionArray[i]) == 0 || m.PositionStdDev[i] == 0 {
			continue
		}

		// Calculate average position for this character in the text
		var posSum float64
		for _, pos := range letterData.PositionArray[i] {
			relPos := float64(pos) / float64(letterData.TotalCount)
			posSum += relPos
		}
		avgPos := posSum / float64(len(letterData.PositionArray[i]))

		// Calculate position z-score
		posZScore := math.Abs(avgPos-m.PositionDistribution[i]) / m.PositionStdDev[i]

		if posZScore > 1.0 {
			var charLabel string
			if i < 10 {
				charLabel = fmt.Sprintf("pos_%d", i)
			} else {
				charLabel = fmt.Sprintf("pos_%s", string(rune('a'+(i-10))))
			}

			zScores[charLabel] = posZScore
			totalScore += posZScore
			significantDeviations++
		}
	}

	// Normalize the score based on number of significant deviations
	if significantDeviations > 0 {
		totalScore /= float64(significantDeviations)
	}

	return totalScore, zScores
}

// IsAnomaly determines if a text is anomalous compared to the model
func (m *TextDistributionModel) IsAnomaly(text string) (bool, float64, map[string]float64) {
	score, anomalies := m.AnomalyScore(text)
	return score > m.AnomalyThreshold, score, anomalies
}

// GetTopAnomalies returns the top n anomalous characters sorted by z-score
func (m *TextDistributionModel) GetTopAnomalies(text string, n int) []string {
	_, anomalies := m.AnomalyScore(text)

	// Convert map to slice for sorting
	type anomalyEntry struct {
		char  string
		score float64
	}

	var entries []anomalyEntry
	for char, score := range anomalies {
		entries = append(entries, anomalyEntry{char, score})
	}

	// Sort by score (descending)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score > entries[j].score
	})

	// Get top n results
	var results []string
	for i := 0; i < min(n, len(entries)); i++ {
		results = append(results, fmt.Sprintf("%s (%.2f)", entries[i].char, entries[i].score))
	}

	return results
}

// ModelSummary returns a string summary of the distribution model
func (m *TextDistributionModel) ModelSummary() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Distribution Model Summary:\n"))
	sb.WriteString(fmt.Sprintf("Based on %d text samples\n", m.SampleCount))
	sb.WriteString(fmt.Sprintf("Anomaly threshold: %.2f\n\n", m.AnomalyThreshold))

	sb.WriteString("Character frequency distribution:\n")
	for i := 0; i < 36; i++ {
		var char string
		if i < 10 {
			char = fmt.Sprintf("%d", i)
		} else {
			char = string(rune('a' + (i - 10)))
		}

		sb.WriteString(fmt.Sprintf("%s: %.4f (±%.4f)\n",
			char, m.CharDistribution[i], m.CharStdDev[i]))
	}

	return sb.String()
}

// CreateDistributionModelWithFitting builds a statistical model with distribution fitting
// from multiple text samples
func CreateDistributionModelWithFitting(textSamples []string, anomalyThreshold float64) (*TextDistributionWithFitting, error) {
	if len(textSamples) == 0 {
		return nil, fmt.Errorf("no text samples provided")
	}

	// Initialize the model
	model := &TextDistributionWithFitting{
		TextDistributionModel: TextDistributionModel{
			SampleCount:      len(textSamples),
			AnomalyThreshold: anomalyThreshold,
		},
	}

	// Collect letter data from all samples
	var allLetterData []*LetterData
	for _, text := range textSamples {
		letterData := AnalyzeLettersFromText(text)
		allLetterData = append(allLetterData, letterData)
	}

	// Process each character
	for i := 0; i < 36; i++ {
		// Collect frequency data for this character across all samples
		var frequencies []float64
		for _, ld := range allLetterData {
			if ld.TotalCount > 0 {
				relFreq := float64(ld.LetterNumberArray[i]) / float64(ld.TotalCount)
				frequencies = append(frequencies, relFreq)
			}
		}

		// Store raw frequency data
		model.CharFrequencyData[i] = frequencies

		// Calculate basic statistics
		if len(frequencies) > 0 {
			model.TextDistributionModel.CharDistribution[i] = stat.Mean(frequencies, nil)
			if len(frequencies) > 1 {
				model.TextDistributionModel.CharStdDev[i] = stat.StdDev(frequencies, nil)
			}
		}

		// Fit distributions if we have enough data
		if len(frequencies) >= 5 {
			model.CharDistributionType[i] = FindBestDistribution(frequencies)
		} else {
			// Default to normal with calculated parameters
			model.CharDistributionType[i] = DistributionParameters{
				Type:   NormalDist,
				Mean:   model.TextDistributionModel.CharDistribution[i],
				StdDev: model.TextDistributionModel.CharStdDev[i],
			}
		}
	}

	// Position distributions (using the original method)
	for i := 0; i < 36; i++ {
		var positions []float64

		// Gather positions for this character across all samples
		for _, ld := range allLetterData {
			if ld.TotalCount > 0 {
				for _, pos := range ld.PositionArray[i] {
					relPos := float64(pos) / float64(ld.TotalCount)
					positions = append(positions, relPos)
				}
			}
		}

		// Calculate position statistics
		if len(positions) > 0 {
			model.TextDistributionModel.PositionDistribution[i] = stat.Mean(positions, nil)
			if len(positions) > 1 {
				model.TextDistributionModel.PositionStdDev[i] = stat.StdDev(positions, nil)
			}
		}
	}

	return model, nil
}

// FindBestDistribution determines which probability distribution best fits the data
// Returns distribution parameters for the best fitting distribution
func FindBestDistribution(data []float64) DistributionParameters {
	if len(data) < 5 {
		// Not enough data for reliable fitting
		mean, std := stat.MeanStdDev(data, nil)
		return DistributionParameters{
			Type:   NormalDist,
			Mean:   mean,
			StdDev: std,
		}
	}

	// Sort the data for CDF calculations
	sortedData := make([]float64, len(data))
	copy(sortedData, data)
	sort.Float64s(sortedData)

	// Test different distributions
	distributions := []struct {
		distType DistributionType
		fitFunc  func([]float64) (DistributionParameters, float64)
	}{
		{NormalDist, fitNormal},
		{GammaDist, fitGamma},
		{BetaDist, fitBeta},
		{ExponentialDist, fitExponential},
		{LogNormalDist, fitLogNormal},
	}

	bestFit := DistributionParameters{Type: NormalDist}
	bestScore := math.Inf(-1)

	for _, dist := range distributions {
		params, score := dist.fitFunc(sortedData)
		if score > bestScore {
			bestScore = score
			bestFit = params
		}
	}

	// Always create empirical distribution as fallback
	empiricalParams := createEmpiricalDistribution(sortedData)

	// If best fit is poor, default to empirical
	if bestScore < 0.8 {
		empiricalParams.GoodnessOfFit = 1.0 // By definition, empirical distribution fits the data perfectly
		return empiricalParams
	}

	bestFit.GoodnessOfFit = bestScore
	return bestFit
}

// fitNormal fits a normal distribution to the data
func fitNormal(data []float64) (DistributionParameters, float64) {
	mean, std := stat.MeanStdDev(data, nil)

	// Create normal distribution
	normal := distuv.Normal{
		Mu:    mean,
		Sigma: std,
	}

	// Evaluate goodness of fit using KS test
	score := goodnessOfFit(data, func(x float64) float64 {
		return normal.CDF(x)
	})

	return DistributionParameters{
		Type:   NormalDist,
		Mean:   mean,
		StdDev: std,
	}, score
}

// fitGamma fits a gamma distribution to the data
func fitGamma(data []float64) (DistributionParameters, float64) {
	// Estimate parameters using method of moments
	mean := stat.Mean(data, nil)
	variance := stat.Variance(data, nil)

	// Gamma parameters: shape (alpha) and rate (beta)
	alpha := mean * mean / variance
	beta := mean / variance

	// If the shape parameter is too small or NaN, this distribution is a poor fit
	if math.IsNaN(alpha) || alpha < 0.1 {
		return DistributionParameters{}, math.Inf(-1)
	}

	// Create gamma distribution
	gamma := distuv.Gamma{
		Alpha: alpha,
		Beta:  beta,
	}

	// Evaluate goodness of fit
	score := goodnessOfFit(data, func(x float64) float64 {
		return gamma.CDF(x)
	})

	return DistributionParameters{
		Type:   GammaDist,
		Shape:  alpha,
		Rate:   beta,
		Mean:   mean,
		StdDev: math.Sqrt(variance),
	}, score
}

// fitBeta fits a beta distribution to the data (assuming data is in range [0,1])
func fitBeta(data []float64) (DistributionParameters, float64) {
	// Check if data is in range [0,1]
	for _, v := range data {
		if v < 0 || v > 1 {
			return DistributionParameters{}, math.Inf(-1)
		}
	}

	mean := stat.Mean(data, nil)
	variance := stat.Variance(data, nil)

	// Beta distribution parameter estimation
	if variance == 0 || mean == 0 || mean == 1 {
		return DistributionParameters{}, math.Inf(-1)
	}

	temp := mean*(1-mean)/variance - 1
	alpha := mean * temp
	beta := (1 - mean) * temp

	// If parameters are invalid, this distribution is a poor fit
	if math.IsNaN(alpha) || math.IsNaN(beta) || alpha <= 0 || beta <= 0 {
		return DistributionParameters{}, math.Inf(-1)
	}

	// Create beta distribution
	betaDist := distuv.Beta{
		Alpha: alpha,
		Beta:  beta,
	}

	// Evaluate goodness of fit
	score := goodnessOfFit(data, func(x float64) float64 {
		return betaDist.CDF(x)
	})

	return DistributionParameters{
		Type:   BetaDist,
		Shape:  alpha,
		Rate:   beta,
		Mean:   mean,
		StdDev: math.Sqrt(variance),
	}, score
}

// fitExponential fits an exponential distribution to the data
func fitExponential(data []float64) (DistributionParameters, float64) {
	// Check if data is non-negative
	for _, v := range data {
		if v < 0 {
			return DistributionParameters{}, math.Inf(-1)
		}
	}

	mean := stat.Mean(data, nil)
	if mean <= 0 {
		return DistributionParameters{}, math.Inf(-1)
	}

	// Lambda is 1/mean for exponential
	lambda := 1.0 / mean

	// Create exponential distribution
	exp := distuv.Exponential{
		Rate: lambda,
	}

	// Evaluate goodness of fit
	score := goodnessOfFit(data, func(x float64) float64 {
		return exp.CDF(x)
	})

	return DistributionParameters{
		Type:   ExponentialDist,
		Rate:   lambda,
		Mean:   mean,
		StdDev: mean, // For exponential, std = mean
	}, score
}

// fitLogNormal fits a log-normal distribution to the data
func fitLogNormal(data []float64) (DistributionParameters, float64) {
	// Check if data is positive
	for _, v := range data {
		if v <= 0 {
			return DistributionParameters{}, math.Inf(-1)
		}
	}

	// Transform to log space
	logData := make([]float64, len(data))
	for i, v := range data {
		logData[i] = math.Log(v)
	}

	// Compute mean and std in log space
	mu, sigma := stat.MeanStdDev(logData, nil)

	// Create lognormal distribution
	lnorm := distuv.LogNormal{
		Mu:    mu,
		Sigma: sigma,
	}

	// Evaluate goodness of fit
	score := goodnessOfFit(data, func(x float64) float64 {
		return lnorm.CDF(x)
	})

	// Calculate actual mean and std in original space
	mean := math.Exp(mu + sigma*sigma/2)
	stdDev := math.Sqrt((math.Exp(sigma*sigma) - 1) * math.Exp(2*mu+sigma*sigma))

	return DistributionParameters{
		Type:   LogNormalDist,
		Mean:   mean,
		StdDev: stdDev,
		Shape:  mu,    // Using Shape to store mu
		Scale:  sigma, // Using Scale to store sigma
	}, score
}

// createEmpiricalDistribution creates an empirical distribution from the data
func createEmpiricalDistribution(data []float64) DistributionParameters {
	n := len(data)

	// Create empirical CDF
	cdf := make([]float64, n)
	for i := range cdf {
		cdf[i] = float64(i+1) / float64(n)
	}

	return DistributionParameters{
		Type:          EmpiricalDist,
		EmpiricalCDF:  cdf,
		Bins:          data,
		Mean:          stat.Mean(data, nil),
		StdDev:        stat.StdDev(data, nil),
		GoodnessOfFit: 1.0, // Empirical distribution has perfect fit by definition
	}
}

// goodnessOfFit calculates the goodness of fit using Kolmogorov-Smirnov test
// Returns a score between 0 and 1, where higher is better
func goodnessOfFit(data []float64, cdf func(float64) float64) float64 {
	n := float64(len(data))
	maxDiff := 0.0

	// Calculate the empirical CDF
	sort.Float64s(data)

	// Find maximum difference between empirical and theoretical CDF
	for i, x := range data {
		// Empirical CDF at point x
		empirical := float64(i+1) / n

		// Theoretical CDF at point x
		theoretical := cdf(x)

		// Calculate difference
		diff := math.Abs(empirical - theoretical)
		if diff > maxDiff {
			maxDiff = diff
		}

		// Also check the previous point to handle discontinuities
		if i > 0 {
			empiricalPrev := float64(i) / n
			diff = math.Abs(empiricalPrev - theoretical)
			if diff > maxDiff {
				maxDiff = diff
			}
		}
	}

	// Convert KS statistic to a score (1 - normalized KS statistic)
	// The critical value for significance level 0.05 is approximately 1.36/sqrt(n)
	criticalValue := 1.36 / math.Sqrt(n)

	if maxDiff > criticalValue {
		// Poor fit - normalize to [0, 0.8)
		return 0.8 * (1.0 - maxDiff/math.Sqrt(n))
	} else {
		// Good fit - normalize to [0.8, 1.0]
		return 0.8 + 0.2*(1.0-maxDiff/criticalValue)
	}
}

// CalculateProbability returns the probability of observing a value
// according to the fitted distribution
func (dp *DistributionParameters) CalculateProbability(value float64) float64 {
	switch dp.Type {
	case NormalDist:
		normal := distuv.Normal{
			Mu:    dp.Mean,
			Sigma: dp.StdDev,
		}
		return normal.Prob(value)

	case GammaDist:
		gamma := distuv.Gamma{
			Alpha: dp.Shape,
			Beta:  dp.Rate,
		}
		return gamma.Prob(value)

	case BetaDist:
		beta := distuv.Beta{
			Alpha: dp.Shape,
			Beta:  dp.Rate,
		}
		return beta.Prob(value)

	case ExponentialDist:
		exp := distuv.Exponential{
			Rate: dp.Rate,
		}
		return exp.Prob(value)

	case LogNormalDist:
		lnorm := distuv.LogNormal{
			Mu:    dp.Shape, // Using Shape to store mu
			Sigma: dp.Scale, // Using Scale to store sigma
		}
		return lnorm.Prob(value)

	case EmpiricalDist:
		// For empirical distribution, use kernel density estimation
		return empiricalProbability(value, dp.Bins)

	default:
		return 0
	}
}

// empiricalProbability estimates probability using kernel density estimation
func empiricalProbability(x float64, data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	// Use Silverman's rule for bandwidth
	n := float64(len(data))
	std := stat.StdDev(data, nil)
	h := 1.06 * std * math.Pow(n, -0.2)

	if h == 0 {
		// If std is 0, use a small bandwidth
		h = 0.01
	}

	// Gaussian kernel density estimation
	sum := 0.0
	for _, xi := range data {
		z := (x - xi) / h
		sum += math.Exp(-0.5 * z * z)
	}

	return sum / (n * h * math.Sqrt(2*math.Pi))
}

// AnomalyScoreWithFitting calculates how different a text is from the fitted distributions
func (m *TextDistributionWithFitting) AnomalyScoreWithFitting(text string) (float64, map[string]float64) {
	// Analyze the text
	letterData := AnalyzeLettersFromText(text)

	// Calculate likelihood scores for each character
	anomalyScores := make(map[string]float64)
	var totalScore float64
	var significantDeviations int

	for i := 0; i < 36; i++ {
		// Skip characters with no distribution data
		if len(m.CharFrequencyData[i]) == 0 {
			continue
		}

		// Calculate relative frequency for this character
		var relFreq float64
		if letterData.TotalCount > 0 {
			relFreq = float64(letterData.LetterNumberArray[i]) / float64(letterData.TotalCount)
		}

		// Calculate probability of observing this frequency
		prob := m.CharDistributionType[i].CalculateProbability(relFreq)

		// Convert to anomaly score (lower probability = higher anomaly)
		// Use negative log probability as anomaly score
		var anomalyScore float64
		if prob > 0 {
			anomalyScore = -math.Log10(prob)
		} else {
			anomalyScore = 10 // Very high anomaly for zero probability
		}

		// Only count significant deviations
		// Using threshold of 2 (prob < 0.01) for significance
		if anomalyScore > 2 {
			var charLabel string
			if i < 10 {
				// It's a number
				charLabel = fmt.Sprintf("%d", i)
			} else {
				// It's a letter
				charLabel = string(rune('a' + (i - 10)))
			}

			anomalyScores[charLabel] = anomalyScore
			totalScore += anomalyScore
			significantDeviations++
		}
	}

	// Normalize the score
	if significantDeviations > 0 {
		totalScore /= float64(significantDeviations)
	}

	return totalScore, anomalyScores
}

// IsAnomalyWithFitting determines if a text is anomalous using fitted distributions
func (m *TextDistributionWithFitting) IsAnomalyWithFitting(text string) (bool, float64, map[string]float64) {
	score, anomalies := m.AnomalyScoreWithFitting(text)
	return score > m.AnomalyThreshold, score, anomalies
}

// GetDistributionSummary returns a string summary of the fitted distributions
func (m *TextDistributionWithFitting) GetDistributionSummary() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Fitted Distribution Model Summary:\n"))
	sb.WriteString(fmt.Sprintf("Based on %d text samples\n", m.SampleCount))
	sb.WriteString(fmt.Sprintf("Anomaly threshold: %.2f\n\n", m.AnomalyThreshold))

	sb.WriteString("Character distribution types:\n")
	for i := 0; i < 36; i++ {
		if len(m.CharFrequencyData[i]) == 0 {
			continue
		}

		var char string
		if i < 10 {
			char = fmt.Sprintf("%d", i)
		} else {
			char = string(rune('a' + (i - 10)))
		}

		dist := m.CharDistributionType[i]

		sb.WriteString(fmt.Sprintf("%s: %s distribution (fit: %.2f)\n",
			char, dist.Type, dist.GoodnessOfFit))

		switch dist.Type {
		case NormalDist:
			sb.WriteString(fmt.Sprintf("   Mean: %.4f, StdDev: %.4f\n", dist.Mean, dist.StdDev))
		case GammaDist:
			sb.WriteString(fmt.Sprintf("   Shape: %.4f, Rate: %.4f\n", dist.Shape, dist.Rate))
		case BetaDist:
			sb.WriteString(fmt.Sprintf("   Alpha: %.4f, Beta: %.4f\n", dist.Shape, dist.Rate))
		case ExponentialDist:
			sb.WriteString(fmt.Sprintf("   Rate: %.4f\n", dist.Rate))
		case LogNormalDist:
			sb.WriteString(fmt.Sprintf("   Mu: %.4f, Sigma: %.4f\n", dist.Shape, dist.Scale))
		case EmpiricalDist:
			sb.WriteString(fmt.Sprintf("   Sample size: %d\n", len(dist.Bins)))
		}
	}

	return sb.String()
}
