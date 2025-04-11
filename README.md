# GoFigure
NOTE: this is work in progress.
A program to compare the similarity of texts in a very simple manner using Go.

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

## Internal Design

- Characters are mapped into a 36-element space: 0–9 for digits and 10–35 for a–z.
- Position-based comparison normalizes indexes relative to total text length.
- Distribution fitting leverages `gonum/stat` for statistical functions.

## Installation

```bash
go get github.com/ML1883/GoFigure
```

## Dependencies

- [Gonum](https://github.com/gonum/gonum) – for statistical functions and distribution fitting