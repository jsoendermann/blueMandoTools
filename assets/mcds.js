// constants
var alphabet = "abcdefghijklmnopqrstuvwxyz".split("");
var nonCharacters = new Array('0', '1', '2', '3', '4',
    '5', '6', '7', '8', '9', 
    '（', '）', '(', ')', "！", "!", "？", "?", "~",
    '「', '」', '《', '》', '【', '】', '…',
    '。', '，', '、', '　', ' ', '；',
    '.', ',', ';', '"', "'", "“", "”", "—", "-", "_",
    '：', ':', "\n", ' ');
nonCharacters = nonCharacters.concat(alphabet, "ABCDEFGHIJKLMNOPQRSTUVWXYZ".split(""));
var easyToReachKeys = new Array('f', 'd', 'k', 'r', 
    'u', 'e', 'i', 's', 'l', 
    'w', 'o', 'j');
var frequentlyUsedChars = "的當了在有为為以就才".split("");
var keysForFrequentlyUsedChars = new Array();
for (var i in alphabet)
{
  var c = alphabet[i];
  if (easyToReachKeys.indexOf(c) == -1)
    keysForFrequentlyUsedChars.push(c);
}
var CLOZE_BEGIN = '<span style="font-weight:600;color:#ff12c7;">';
var CLOZE_END = '</span>';
var CLOZE_SYMBOL = CLOZE_BEGIN+'％'+CLOZE_END;


// dicts
var dict = {};
var reverseDict = {};
var pos = 0;


function generateDict(characters)
{
  dict = {};
  reverseDict = {};
  pos = 0;

  var iEtrk = 0, iAlphabet = 0;
  var iFuc = 0;

  for (var i in characters)
  {
    var v = characters[i];

    if (!dict[v])
    {
      if (frequentlyUsedChars.indexOf(v) != -1)
      {
        dict[v] = keysForFrequentlyUsedChars[iFuc];
        reverseDict[keysForFrequentlyUsedChars[iFuc]] = v;
        iFuc++;
      }
      else
      {
        var c = easyToReachKeys[iEtrk]+alphabet[iAlphabet];
        dict[v] = c;
        reverseDict[c] = v;
        iAlphabet++;
      }
    }
    if (iAlphabet >= alphabet.length)
    {
      iAlphabet = 0;
      iEtrk++;
    }
  }

  return;
}

function generateMCD(clozedChar)
{
  var text = $("#mcds-input").val();
  var notes = $("#mcds-notes").val();
  notes = notes.replace(/\n/g,'<br />');
  notes = notes.replace(/\t/g,'&nbsp;&nbsp;&nbsp;&nbsp;');
  text = text.replace(/\n/g,'<br />');
  text = text.replace(/\t/g,'&nbsp;&nbsp;&nbsp;&nbsp;');
  var regex = new RegExp(clozedChar, "g");
  var front = text.replace(regex, CLOZE_SYMBOL);
  var back = text.replace(regex, CLOZE_BEGIN+clozedChar+CLOZE_END);

  return front+"\t"+back+"\t"+notes;
}

function generateHTML()
{
  var text = $("#mcds-input").val();
  text = text.replace(/\n/g,"☃");
  var output = "";
  var textSplit = text.split("");

  for (var i in textSplit)
  {
    var v = textSplit[i];

    if (v == "☃")
      output += "<br />";

    else if (dict[v])
      output += "<ruby><rb>"+v+'</rb><rt class="rt-'+dict[v]+'" id="rt-'+i+'">'+dict[v]+"</rt></ruby>";
    else
      output += v;
  }
  return output;
}

function processInput(rawInput)
{
  var allCharacters = eliminateDuplicates(rawInput.split(""));

  // remove non chinese characters from array
  var chineseCharacters = new Array();

  for (var i in allCharacters)
  {
    var v = allCharacters[i];
    if (nonCharacters.indexOf(v) == -1)
      chineseCharacters.push(v);
  }

  generateDict(chineseCharacters);

  $("#mcds-display").html(generateHTML());
}


function eliminateDuplicates(arr) {
  var i,
      len=arr.length,
      out=[],
      obj={};

  for (i=0;i<len;i++) {
    obj[arr[i]]=0;
  }
  for (i in obj) {
    out.push(i);
  }
  return out;
}

function processCodeInputKeyup()
{
  var code = $("#mcds-code-input").val();

  if (reverseDict[code])
  {
    //var original = "<ruby><rb>"+reverseDict[code]+"</rb><rt>"+code+"</rt></ruby>";
    //var regex = new RegExp(original, "g");
    //var replaceBy = '<ruby><rb class="already-clozed-rb">'+reverseDict[code]+'</rb><rt class="already-clozed-rt">'+code+"</rt></ruby>"

    // $("#mcds-display").html($("#mcds-display").html().replace(regex, replaceBy));
    var rts = $(".rt-"+code);

    rts.fadeTo("fast", 0);
    
    $("#mcds-clozed-chars").val($("#mcds-clozed-chars").val() + reverseDict[code] + " ");


    $("#mcds-code-input").val("");
    
    /*var mcdsOutputVal = $("#mcds-output").val();
    $("#mcds-output").val(mcdsOutputVal+(mcdsOutputVal==""?"":"\n")+generateMCD(reverseDict[code]));
    var mcdsOutputTa = $("#mcds-output");
    mcdsOutputTa.scrollTop(mcdsOutputTa[0].scrollHeight - mcdsOutputTa.height());*/

    // #####

    /*$("#mcds-notes").val($("#mcds-notes").val()+rts.length);
      for (var i in rts.length) {

    //var v = rts[i];
    //$("#mcds-notes").val($("#mcds-notes").val()+v);
    $("#mcds-notes").val($("#mcds-notes").val()+i);
    }*/
  }

}


$(function(){
  if (! ($("#mcds-input").val() == undefined)) {
    console.log("Executing mcds.js functions")

    $("#mcds-input").keyup(function() {
      processInput($("#mcds-input").val());
    })
    $("#mcds-code-input").keyup(function() {
      processCodeInputKeyup();
    })
    $("#mcds-output").click(function() {
      $("#mcds-output").focus();
      $("#mcds-output").select();
    })

    processInput($("#mcds-input").val());
  }
});


