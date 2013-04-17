# Document ready
$( ->
  # Button: Lookup Words
  $('#sc-lookup').on('click', ->
    scLookupClicked()
  )

  # Select all on focus
  selectAllOnFocus('#sc-output')
)

# Lookup button event handler
scLookupClicked = ->
  # get words
  sentences = $('#sc-sentences').val().split("\n")

  # get colors
  tones = getColors()

  for sentence in sentences
    # make ajax request to server
    $.ajax({url: "/sentences/lookup/#{sentence}", async: true, dataType: 'json', data: tones}).success( (response) ->
      # if there was no error, add the response to #sc-output
      # TODO deal with error
      if response["error"] == 'nil'
        textAreaAddLineAndScroll '#sc-output', response['csv']
        #$('#debug').html(response['csv'])
      else
        console.log response["error"]
    )

