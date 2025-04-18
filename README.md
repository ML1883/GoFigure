# GoFigure
  This is a Go program/library to compare the similarity of texts using letter frequency and letter position. The `letteranalyzer` module compares two texts based on this data and comes up with an exact number how similair the texts are. The `letterdistribution` module builds upon the data structure from `letteranalyzer` and fits statistical distributions to the frequency and position data across a group of texts. This fitted distribution can then be used to determine whether a new text follows the same pattern as the training texts, effectively serving as a model for comparison./ You might use this to e.g. check if a text is an anomaly compared to others. I am not sure how well this method works compared to different techniques of comparing text similairity, but the idea of comparing texts using exact/analytical measures (i.e. frequency counting and such) as opposed to statistical interested me, and that is why I created this module. I think the crude method this library uses has it's up and downsides. It is up to the end user to determine what carries more weight.

## Features

- **Character Frequency Analysis**  
  Parses input text and counts occurrences and positions of all alphanumeric characters.

- **Vector Similarity Metrics**
  - Cosine Similarity
  - Jaccard Index
  - Position Difference Score (a custom measure to handle the odd shapes these might take)

- **Statistical Distribution Fitting & Scoring**  
  Builds statistical models from multiple text samples, estimating mean, standard deviation, and fitting probability distributions (normal, gamma, beta, etc.) to the frequency and position data of the characters found in the texts. Automatically selects best-fit distributions using statistical metrics for each character.

## Installation

```bash
go get github.com/ML1883/GoFigure
```

## Internal Design

- Characters are mapped into a 36-element space: 0–9 for digits and 10–35 for a–z.
- Position-based comparison normalizes indexes relative to total text length.
- Distribution fitting leverages `gonum/stat` for statistical functions.

## Command-Line Interface

The `cmd/test` directory contains a CLI application that makes the package's functionality available through command-line arguments:

### Comparison Mode
```bash
# Compare two files
./main -compare -file -text1=file1.txt -text2=file2.txt

# Interactive/direct input text comparison
./main -compare
```

### Distribution Mode

```bash
# Create a distribution model from training texts
./main -distribution -create-model -folder=./training_texts -model-file=model.gob

# Check text against an existing model
./main -distribution -use-model -model-file=model.gob -check-text=sample.txt

# Interactive/direct input analysis with existing model
./main -distribution -use-model -model-file=model.gob
```

### Additional Options

- `-output`: Show detailed vectors and statistical arrays
- `-threshold=2.0`: Adjust anomaly detection sensitivity (higher = more strict)
- `-fit-threshold=0.8`: Control distribution fitting (higher = more empirical)
- `-help`: Display help information

## Known problems/TODO
- Extremely similair model training texts causing distribution shapes to go to infinity.
- calculatedProb variations of the function AnomalyScore can be NaN.

## Dependencies

- [Gonum](https://github.com/gonum/gonum) – for statistical functions and distribution fitting