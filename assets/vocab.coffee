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
      vcOutputAddLine response["response"]
    )

vcOutputAddLine = (line) ->
  console.log 'vcOutputAddLine'
  console.log line

  $('#vc-output').val($('#vc-output').val() + line + "\n")

  mcdsOutputTa = $("#vc-output")
  mcdsOutputTa.scrollTop(mcdsOutputTa[0].scrollHeight - mcdsOutputTa.height())
