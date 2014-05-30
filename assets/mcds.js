/* global $ */

(function ($) {
    "use strict";

    var alphabet, nonChineseCharacters, firstLettersInCodes, charToCodeDict, codeToCharDict;

    // constants
    alphabet = "abcdefghijklmnopqrstuvwxyz";

    nonChineseCharacters = '0123456789（）()！!？?~～"「」《》【】…\'。，、　 ；.,;"“”—-_：:\n", \t'.split();
    nonChineseCharacters += alphabet.split();
    nonChineseCharacters += alphabet.toUpperCase().split('');

    firstLettersInCodes = 'fjdkslaurieowpq'.split('');

    // objects shared among the functions
    charToCodeDict = {};
    codeToCharDict = {};


    $(function(){
        $("#mcds-input").keyup(function() {
            processInputText();
        });

        $("#mcds-code-input").keyup(function() {
            processCodeInputKeyup();
        });

        processInputText($("#mcds-input").val());
    });

    function processInputText() {
        var text, uniqueCharacters, chineseCharacters, i, character;

        text = $("#mcds-input").val();
        
        uniqueCharacters = eliminateDuplicates(text.split(""));

        // remove non chinese characters from array
        chineseCharacters = [];
        for (i = 0; i < uniqueCharacters.length; i++) {
            character = uniqueCharacters[i];
            if (nonChineseCharacters.indexOf(character) === -1) {
                chineseCharacters.push(character);
            }
        }

        generateDict(chineseCharacters);
        updateTextDisplay();
        $("#mcds-clozed-chars").val("");
    }

    function processCodeInputKeyup() {
        var code, codeRubyClass;
        code = $("#mcds-code-input").val();

        if (typeof codeToCharDict[code] !== 'undefined') {
            codeRubyClass = $(".rt-" + code);
            codeRubyClass.fadeTo("fast", 0);

            $("#mcds-clozed-chars").val($("#mcds-clozed-chars").val() + codeToCharDict[code] + " ");
            $("#mcds-code-input").val("");
        }
    }


    function generateDict(characters) {
        var possibleCodes, possibleCodesIndex, i, character, code;

        charToCodeDict = {};
        codeToCharDict = {};

        possibleCodes = cartesianProduct(firstLettersInCodes, alphabet.split(''));
        possibleCodesIndex = 0;

        for (i = 0; i < characters.length; i++) {
            character = characters[i];

            if (typeof charToCodeDict[character] === 'undefined') {
                code = possibleCodes[possibleCodesIndex];

                charToCodeDict[character] = code;
                codeToCharDict[code] = character;
                possibleCodesIndex++;
            }
        }
    }

    function updateTextDisplay() {
        var text, output, i, character;

        text = $("#mcds-input").val();

        output = "";

        for (i = 0; i < text.length; i++) {
            character = text[i];

            if (typeof charToCodeDict[character] === 'undefined') {
                output += character;
            } else {
                output += constructRubyHtmlForChar(character);
            }
        }

        output = output.replace(/\n/g, '<br />');

        $("#mcds-display").html(output);
    }

    function constructRubyHtmlForChar(character) {
        return '<ruby><rb>' + character + '</rb>' +
            '<rt class="rt-' + charToCodeDict[character] + '">' +
            charToCodeDict[character] + "</rt></ruby>";
    }

    function eliminateDuplicates(array) {
        var uniqueArray = [];

        $.each(array, function(i, element) {
            if($.inArray(element, uniqueArray) === -1) {
                uniqueArray.push(element);
            }
        });

        return uniqueArray;
    }

    function cartesianProduct(array1, array2) {
        var product = [];

        $.each(array1, function(i, element1) {
            $.each(array2, function(j, element2) {
                product.push(element1 + element2);
            });
        });

        return product;
    }
})($);
