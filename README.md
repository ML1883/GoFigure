# GoFigure
NOTE: this is work in progress.  

This is a program to compare the similarity of texts in Go using letter frequency and position counting. I am not sure how well this method works compared to different techniques of comparing text similairity, but the idea of comparing texts using exact/analytical measures (i.e. frequency counting and such) as opposed to statistical interested me, and that is why I created this module. I also wanted to try out the programming language Go for the first time.

## Overview

The `analyzer` package processes and compares text based on the distribution and position of characters (letters a–z and digits 0–9). It’s designed for use cases like anomaly detection, text similarity, and distribution modeling. The `parses` package provides basic parsing and file operations.

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

The `cmd` directory contains a CLI application that makes the package's functionality available through command-line arguments:

### Comparison Mode

```bash
# Compare two files
./program -compare -file -text1=file1.txt -text2=file2.txt

# Interactive text comparison
./program -compare
```

### Distribution Mode

```bash
# Create a distribution model from training texts
./program -distribution -create-model -folder=./training_texts -model-file=model.gob

# Check text against an existing model
./program -distribution -use-model -model-file=model.gob -check-text=sample.txt

# Interactive analysis with existing model
./program -distribution -use-model -model-file=model.gob
```

### Additional Options

- `-output`: Show detailed vectors and statistical arrays
- `-threshold=2.0`: Adjust anomaly detection sensitivity (higher = more strict)
- `-fit-threshold=0.8`: Control distribution fitting (higher = more empirical)
- `-help`: Display help information

## Dependencies

- [Gonum](https://github.com/gonum/gonum) – for statistical functions and distribution fitting