$(document).ready( ->
  # Button: Lookup Words
  $('#mvc-lookup').on('click', ->
    mvcLookupClicked()
  )

  selectAllOnFocus('#mvc-output')
  selectAllOnFocus('#mvc-not-found')
)


# Lookup button event handler
mvcLookupClicked = ->
  words = $('#mvc-words').val().split("\n")

  tones = getColors()

  for word in words
    # make ajax request to server
    $.ajax({url: "/moe-vocab/lookup/#{word}", async: true, dataType: 'json', data: tones}).success( (response) ->
      if response["error"] == 'nil'
        textAreaAddLineAndScroll '#mvc-output', response['csv']
      else
        textAreaAddLineAndScroll '#mvc-not-found', response['word']

    )

