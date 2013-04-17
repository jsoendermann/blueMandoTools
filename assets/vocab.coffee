# Document ready
$( ->
  # Button: Lookup Words
  $('#vc-lookup').on('click', ->
    vcLookupClicked()
  )
)

# Lookup button event handler
vcLookupClicked = ->
  # get words
  words = $('#vc-words').val().split("\n")

  for word in words
    # make ajax request to server
    $.ajax({url: "/vocab/lookup/#{word}", async: true, dataType: 'json'}).success( (response) ->
      # if there was no error, add the response to #vc-output...
      if response["error"] == 'nil'
        textAreaAddLineAndScroll '#vc-output', response['response']
      # ...otherwise add the word to #vc-not-found
      else
        textAreaAddLineAndScroll '#vc-not-found', response['word']

    )

# this function adds a line of text to a text area and scrolls down so it's visible
textAreaAddLineAndScroll = (textAreaId, line) ->
  ta = $(textAreaId)

  ta.val(ta.val() + line + '\n')
  ta.scrollTop(ta[0].scrollHeight - ta.height())

