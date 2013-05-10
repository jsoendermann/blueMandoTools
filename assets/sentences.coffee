$(document).ready( ->
  # Button: Lookup Words
  $('#sc-lookup').on('click', ->
    scLookupClicked()
  )

  selectAllOnFocus('#sc-output')
)

# Lookup button event handler
scLookupClicked = ->
  sentences = $('#sc-sentences').val().split("\n")

  tones = getColors()

  for sentence in sentences
    # make ajax request to server
    $.ajax({url: "/sentences/lookup/#{sentence}", async: true, dataType: 'json', data: tones}).success( (response) ->
      # TODO handle error
      if response["error"] == 'nil'
        textAreaAddLineAndScroll '#sc-output', response['csv']
      else
        console.log response["error"]
    )

