import theme from '../../../theme/theme';

export const globalStyles = {
  '@global': {
    html: {
      height: '100%'
    },
    body: {
      height: '100%'
    },
    '#root': {
      height: '100%',
      display: 'flex',
      flexDirection: 'column'
    },
    a: {
      textDecoration: 'underline',
      wordBreak: 'break-word',
      color: theme.palette.primary.main
    },
    img: {
      maxWidth: '100%'
    },
    iframe: {
      width: '100%',
      height: '500px',
      border: 0
    },
    '.pointer': {
      cursor: 'pointer'
    },
    '.ellipsis': {
      whiteSpace: 'nowrap',
      overflow: 'hidden',
      textOverflow: 'ellipsis'
    },
    '.noDecoration': {
      '&:hover': {
        textDecoration: 'none'
      }
    },
    '.truncate': {
      whiteSpace: 'nowrap',
      overflow: 'hidden',
      textOverflow: 'ellipsis'
    },
    '.navigationItemIcon': {
      width: '26px',
      height: 'auto',
      marginRight: '12px',
      imageRendering: '-webkit-optimize-contrast'
    }
  }
};
