$(document).ready( ->
    # Button: Lookup Words
    $('#mcds-lookup').on('click', ->
        mcdsLookupClicked()
    )

    selectAllOnFocus('#mcds-dict-output')
)

# Lookup button event handler
mcdsLookupClicked = ->
    text = $('#mcds-input').val().replace(/\//g, '@SLASH@')
    notes = $('#mcds-notes').val().replace(/\//g, '@SLASH@')

    chars = $('#mcds-clozed-chars').val()

    tones = getColors()

    $.ajax({url: "/mcds/lookup/#{encodeURIComponent(text)}?chars=#{encodeURIComponent(chars)}&notes=#{encodeURIComponent(notes)}", async: true, dataType: 'json', data: tones}).success( (response) ->
        if response['error'] == 'nil'
            textAreaAddLineAndScroll '#mcds-dict-output', response['result']
    )
