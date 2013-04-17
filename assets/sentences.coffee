# Document ready
$( ->
  # Button: Lookup Words
  $('#sc-lookup').on('click', ->
    scLookupClicked()
  )

  # Select all on focus
  selectAllOnFocus('#vc-output')
  selectAllOnFocus('#vc-not-found')
)

# Lookup button event handler
scLookupClicked = ->
  # get words
  sentences = $('#sc-sentences').val().split("\n")

  # get colors
  tone0 = $('input[name="sc-tone-0"]').val()
  tone1 = $('input[name="sc-tone-1"]').val()
  tone2 = $('input[name="sc-tone-2"]').val()
  tone3 = $('input[name="sc-tone-3"]').val()
  tone4 = $('input[name="sc-tone-4"]').val()

  for sentence in sentences
    # make ajax request to server
    $.ajax({url: "/sentences/lookup/#{sentence}", async: true, dataType: 'json',data: {tone0: tone0, tone1: tone1, tone2: tone2, tone3: tone3, tone4: tone4}}).success( (response) ->
      # if there was no error, add the response to #sc-output
      # TODO deal with error
      if response["error"] == 'nil'
        textAreaAddLineAndScroll '#sc-output', response['csv']
        $('#debug').html(response['csv'])
      else
        console.log response["error"]
    )

