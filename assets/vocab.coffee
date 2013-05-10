# Document ready
$( ->
  # Button: Lookup Words
  $('#vc-lookup').on('click', ->
    vcLookupClicked()
  )

  selectAllOnFocus('#vc-output')
  selectAllOnFocus('#vc-not-found')
)


# Lookup button event handler
vcLookupClicked = ->
  words = $('#vc-words').val().split("\n")

  tones = getColors()

  for word in words
    # make ajax request to server
    $.ajax({url: "/vocab/lookup/#{word}", async: true, dataType: 'json', data: tones}).success( (response) ->
      if response["error"] == 'nil'
        textAreaAddLineAndScroll '#vc-output', response['csv']
      else
        textAreaAddLineAndScroll '#vc-not-found', response['word']

    )

