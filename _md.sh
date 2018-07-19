function _mt() {
  typeset -r trash="${HOME}/.Trash"

  _arguments \
    '-r[restore]: :->trash' \
    '-d[delete]: :->trash' \
    '-l[list]: :->list' \
    '-s[size]: :->none' \
    '*: :->files'

  case "${state}" in
    list )
      _arguments '-day[? day ago]: :->none'
    ;;
    trash )
      _values 'files in trash can' $(command ls -Ar ${trash})
    ;;
    files )
      _files
    ;;
  esac

}
compdef _mt mt
