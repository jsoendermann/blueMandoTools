$(document).ready( ->
  # Button: Lookup Words
  $('#mcds-lookup').on('click', ->
    mcdsLookupClicked()
  )

  selectAllOnFocus('#mcds-dict-output')
)

# Lookup button event handler
mcdsLookupClicked = ->
  mcds = $('#mcds-output').val().split("\n")

  tones = getColors()

  for mcd in mcds
    #FIXME find a better solution 
    mcd = mcd.replace(/\//g, '@SLASH@');
    # make ajax request to server
    $.ajax({url: "/mcds/lookup/#{encodeURIComponent(mcd)}", async: true, dataType: 'json', data: tones}).success( (response) ->
      # TODO handle error
      if response["error"] == 'nil'
        textAreaAddLineAndScroll '#mcds-dict-output', response['csv']
      else
        console.log response["error"]
    )

