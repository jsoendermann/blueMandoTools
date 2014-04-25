package main

// Paths on the server
const (
	moeVocabPath   = "/moe-vocab/"
	vocabPath      = "/vocab/"
	sentencesPath  = "/sentences/"
	mcdsPath       = "/mcds/"

	moeVocabLookupPath  = "/moe-vocab/lookup/"
	vocabLookupPath     = "/vocab/lookup/"
	sentencesLookupPath = "/sentences/lookup/"
	mcdsLookupPath      = "/mcds/lookup/"

	settingsPath  = "/settings/"

	assetsPath = "/assets/"
)

const layoutFile = "layout.html"

// mcd related constants
const (
	clozeBegin = `<span style="font-weight:600;color:#ff12c7;">`
    clozeChar  = `％`
	clozeEnd   = `</span>`
)
