import * as createPalette from '@material-ui/core/styles/createPalette';

declare module '@material-ui/core/styles/createPalette' {
  interface PaletteOptions {
    custom?: CustomThemeColors;
    boxShadows?: CustomBoxShadows;
  }

  interface Palette {
    custom: CustomThemeColors;
    boxShadows: CustomBoxShadows;
  }
}

interface CustomThemeColors {
  mainGray: string;
  dotRed: string;
  lightGray: string;
  white: string;
  darkGray: string;
  transparentBlack: string;
}

interface CustomBoxShadows {
  main: string;
  darker: string;
}
