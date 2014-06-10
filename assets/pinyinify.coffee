$(document).ready( ->
  # Button: Lookup Words
  $('#pinyinify-lookup').on('click', ->
    pinyinifyLookupClicked()
  )

  selectAllOnFocus('#pinyinify-output')
)


# Lookup button event handler
pinyinifyLookupClicked = ->
  text = $('#pinyinify-text').val()

  $.ajax({url: "/pinyinify/lookup/#{text}", async: true, dataType: 'json'}).success( (response) ->
      if response["error"] == 'nil'
          $('#pinyinify-output').val(response['result'])
  )
