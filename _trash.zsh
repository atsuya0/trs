function _trash() {
  function sub_commands() {
    _values 'Commands' \
      'move' \
      'restore' \
      'list' \
      'size' \
      'delete' \
      'auto-delete'
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
        ;;
        (list)
          _arguments \
            '(-d --days)'{-d,--days}'[Display files that are not past (n) days since they were discarded in the trash can.]' \
            '(-s --size)'{-s,--size}'[Display files with size greater than (n).]' \
            '(-r --reverse)'{-r,--reverse}'[display in reverse order]'
        ;;
        (size)
        ;;
        (delete)
        ;;
        (auto-delete)
        ;;
      esac
    ;;
  esac
}
compdef _trash trash
