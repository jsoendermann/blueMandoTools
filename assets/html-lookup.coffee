$(document).ready( ->
	# Button: Lookup Words
	$("#html-input-word").keyup( (event) ->
		if(event.keyCode == 13)
			$("#html-input-lookup").click()
	)

	$('#html-input-lookup').on('click', ->
		htmlInputWordChanged()
	)

	selectAllOnFocus('#html-input-cedict')
	selectAllOnFocus('#html-input-moedict')
)


htmlInputWordChanged = ->
	word = $('#html-input-word').val()
	tones = getColors()

	$.ajax({url: "/vocab/lookup/#{word}", async: true, dataType: 'json', data: tones}).success( (response) ->
		if response["error"] == 'nil'
			if response["csv"].indexOf($('#html-input-word').val()) == 0
				$('#html-input-cedict').val(response.csv)
	)

	$.ajax({url: "/moe-vocab/lookup/#{word}", async: true, dataType: 'json', data: tones}).success( (response) ->
		if response["error"] == 'nil'
			if response["csv"].indexOf($('#html-input-word').val()) == 0
				$('#html-input-moedict').val(response.csv)
	)
