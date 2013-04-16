# Document ready
$( ->
  # Button: Lookup Words
  $('#vc-lookup').on('click', ->
    vcLookupClicked()
  )
)



vcLookupClicked = ->
  console.log 'vcLookupClicked'

  words = $('#vc-words').val().split("\n")

  for word in words
    console.log word
    $.ajax({url: "/vocab/lookup/#{word}", async: true, dataType: 'json'}).success( (response) ->
      console.log response
      if response["error"] == 'nil'
        textAreaAddLineAndScroll '#vc-output', response['response']
      else
        textAreaAddLineAndScroll '#vc-not-found', response['word']

    )

textAreaAddLineAndScroll = (textAreaId, line) ->
  ta = $(textAreaId)

  ta.val(ta.val() + line + '\n')
  ta.scrollTop(ta[0].scrollHeight - ta.height())

