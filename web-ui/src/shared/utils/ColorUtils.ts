import theme from '../../theme/theme';

export default class ColorUtils {
  // Returns the corresponding color codes for the file icons
  static getColorCode(extension: string): ColorCode {
    switch (extension) {
      case 'go':
      case 'png':
      case 'jpeg':
      case 'jpg':
      case 'doc':
      case 'docx':
        return {
          iconColor: '#47b8f7',
          backgroundGradient: theme.palette.workspaceGradients.lightBlue
        };
      case 'zip':
      case '7zip':
      case 'rar':
      case 'xls':
      case 'xlsx':
        return {
          iconColor: '#38c1b9',
          backgroundGradient: theme.palette.workspaceGradients.lightGreen
        };
      case 'ppt':
      case 'pptx':
        return {
          iconColor: '#f7974a',
          backgroundGradient: theme.palette.workspaceGradients.lightYellow
        };
      case 'pdf':
        return {
          iconColor: '#ec5757',
          backgroundGradient: theme.palette.workspaceGradients.lightPink
        };
      default:
        return {
          iconColor: '#998b79',
          backgroundGradient: theme.palette.workspaceGradients.lightBrown
        };
    }
  }
}

export interface ColorCode {
  iconColor: string;
  backgroundGradient: string;
}
