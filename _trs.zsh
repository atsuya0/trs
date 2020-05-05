function _trs() {
  function sub_commands() {
    _values 'Commands' \
      'move' \
      'restore' \
      'list' \
      'size' \
      'remove' \
      'auto-remove'
  }

  _arguments -C \
    '(-h --help)'{-h,--help}'[show help]' \
    '1: :sub_commands' \
    '*:: :->args'

  case ${state} in
    (args)
      case ${words[1]} in
        (move)
          _files
        ;;
        (restore)
          _arguments \
            '(-a --all)'{-a,--all}'[Target all the files]'
        ;;
        (list)
          _arguments \
            '(-d --days)'{-d,--days}'[Show the file names that are not past (n) days since they were discarded in the trash can.]' \
            '(-s --size)'{-s,--size}'[Show the file names with size greater than (n).]' \
            '(-r --reverse)'{-r,--reverse}'[Show the file names in reverse order]' \
            '(-p --path)'{-p,--path}'[Show the file paths]'
        ;;
        (size)
        ;;
        (remove)
        ;;
        (auto-remove)
        ;;
      esac
    ;;
  esac
}
compdef _trs trs
