# Document ready
$( ->
  # Button: Lookup Words
  $('#vc-lookup').on('click', ->
    vcLookupClicked()
  )

  # Select all on focus
  selectAllOnFocus('#vc-output')
  selectAllOnFocus('#vc-not-found')
)


# Lookup button event handler
vcLookupClicked = ->
  # get words
  words = $('#vc-words').val().split("\n")

  # get colors
  tones = getColors()

  for word in words
    # make ajax request to server
    $.ajax({url: "/vocab/lookup/#{word}", async: true, dataType: 'json', data: tones}).success( (response) ->
      # if there was no error, add the response to #vc-output...
      if response["error"] == 'nil'
        textAreaAddLineAndScroll '#vc-output', response['csv']
      # ...otherwise add the word to #vc-not-found
      else
        textAreaAddLineAndScroll '#vc-not-found', response['word']

    )

