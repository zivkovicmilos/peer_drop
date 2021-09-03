import { createMuiTheme, Theme } from '@material-ui/core/styles';

// Create a theme instance.
const theme: Theme = createMuiTheme({
  palette: {
    primary: {
      main: '#000000',
      light: '#777',
      contrastText: '#FFF'
    },
    text: {
      primary: '#101010',
      secondary: '#9A9FA5'
    },
    background: {
      default: '#FFF'
    },
    error: {
      main: '#f53b56'
    },
    custom: {
      mainGray: '#F0F0F0',
      dotRed: '#D12D2D',
      lightGray: '#EFF2F5',
      white: '#FFFFFF',
      darkGray: '#F4F4F4',
      transparentBlack: 'rgba(0, 0, 0, 0.5)'
    },
    boxShadows: {
      main: '1px 3px 6px 0px rgba(128,142,155,0.1)',
      darker: '1px 3px 6px 0px rgba(128,142,155,0.3)'
    }
  },
  typography: {
    fontFamily: `"Montserrat", sans-serif`,
    h1: {
      fontSize: '2.25rem'
    },
    body1: {
      fontSize: '1rem',
      textAlign: 'justify',
      textJustify: 'inter-word'
    }
  }
});

theme.overrides = {
  MuiButton: {
    root: {
      textTransform: 'initial'
    },
    disabled: {
      opacity: 0.5
    }
  },
  MuiInputBase: {
    input: {
      fontFamily: 'Montserrat',
      fontSize: theme.typography.pxToRem(16),
      fontWeight: 'normal'
    }
  },
  MuiDivider: {
    root: {
      backgroundColor: theme.palette.custom.lightGray
    }
  },
  MuiTabs: {
    scrollButtons: {
      width: '30px'
    }
  }
};

export default theme;
