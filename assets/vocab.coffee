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
  tone0 = $('input[name="vc-tone-0"]').val()
  tone1 = $('input[name="vc-tone-1"]').val()
  tone2 = $('input[name="vc-tone-2"]').val()
  tone3 = $('input[name="vc-tone-3"]').val()
  tone4 = $('input[name="vc-tone-4"]').val()

  for word in words
    # make ajax request to server
    $.ajax({url: "/vocab/lookup/#{word}", async: true, dataType: 'json',data: {tone0: tone0, tone1: tone1, tone2: tone2, tone3: tone3, tone4: tone4}}).success( (response) ->
      # if there was no error, add the response to #vc-output...
      if response["error"] == 'nil'
        textAreaAddLineAndScroll '#vc-output', response['csv']
      # ...otherwise add the word to #vc-not-found
      else
        textAreaAddLineAndScroll '#vc-not-found', response['word']

    )

