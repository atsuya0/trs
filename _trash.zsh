function _trash() {
  [[ -n ${TRASH_PATH} ]] && typeset -r trash=${TRASH_PATH} || typeset -r trash="${HOME}/.Trash"
  local ret=1

  function sub_commands() {
    local -a _c

    _c=(
      'move' \
      'restore' \
      'list' \
      'size' \
      'delete'
    )

    _describe -t commands Commands _c
  }

  _arguments -C \
    '(-h --help)'{-h,--help}'[show help]' \
    '1: :sub_commands' \
    '*:: :->args' \
    && ret=0

  case ${state} in
    (args)
      case ${words[1]} in
        (move)
          _files
        ;;
        (restore)
        ;;
        (list)
          _arguments \
            '(-d --days)'{-d,--days}'[How many days ago]' \
            '(-r --reverse)'{-r,--reverse}'[display in reverse order]'
        ;;
        (size)
        ;;
        (delete)
        ;;
      esac
    ;;
  esac

  return ret
}
compdef _trash trash
