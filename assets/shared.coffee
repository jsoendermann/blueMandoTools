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

# From http://www.w3schools.com/js/js_cookies.asp
@getCookie = (c_name) ->
  c_value = document.cookie
  c_start = c_value.indexOf(" " + c_name + "=")
  c_start = c_value.indexOf(c_name + "=")  if c_start is -1
  if c_start is -1
    c_value = null
  else
    c_start = c_value.indexOf("=", c_start) + 1
    c_end = c_value.indexOf(";", c_start)
    c_end = c_value.length  if c_end is -1
    c_value = unescape(c_value.substring(c_start, c_end))
  c_value

@getColors = ->
  tones = {}
  for n in [0..4]
    if getCookie("tone"+n) == "undefined" or getCookie("tone"+n) == null
      switch n
        when 0 then tones['tone'+n] = "#000000"
        when 1 then tones['tone'+n] = "#00ac00"
        when 2 then tones['tone'+n] = "#021bff"
        when 3 then tones['tone'+n] = "#996633"
        when 4 then tones['tone'+n] = "#ff0000"
    else
      tones['tone'+n] = unescape(getCookie('tone'+n))
  tones

