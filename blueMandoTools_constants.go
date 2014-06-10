package main

// Paths on the server
// TODO Turn this into something that the main function can loop over
const (
    homePath       = "/"
	moeVocabPath   = "/moe-vocab/"
	vocabPath      = "/vocab/"
	sentencesPath  = "/sentences/"
	mcdsPath       = "/mcds/"
    pinyinifyPath  = "/pinyinify/"

	moeVocabLookupPath  = "/moe-vocab/lookup/"
	vocabLookupPath     = "/vocab/lookup/"
	sentencesLookupPath = "/sentences/lookup/"
	mcdsLookupPath      = "/mcds/lookup/"
    pinyinifyLookupPath       = "/pinyinify/lookup/"

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
