# FIXME it would be better to combine all the coffee script into one file before
# compiling instead of making these functions global

@selectAllOnFocus = (ta) ->
  $(ta).focus( ->
    $this = $(this)
    $this.select()

    $this.mouseup( ->
      $this.unbind("mouseup")
      return false
    )
  )

# this function adds a line of text to a text area and scrolls down so the new line
# is visible
@textAreaAddLineAndScroll = (textAreaId, line) ->
  ta = $(textAreaId)

  ta.val(ta.val() + line + '\n')
  ta.scrollTop(ta[0].scrollHeight - ta.height())

@getColors = ->
  tones = {}
  for n in [0..4]
    tones['tone'+n] = $('input[name="vsc-tone-'+n+'"]').val()
  return tones
