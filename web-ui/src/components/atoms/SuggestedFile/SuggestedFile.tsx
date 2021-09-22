import { Box, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import clsx from 'clsx';
import { FC } from 'react';
import FileIcon, { IconStyle } from 'react-fileicons';
import ColorUtils from '../../../shared/utils/ColorUtils';
import theme from '../../../theme/theme';
import { ISuggestedFileProps } from './suggestedFile.types';

const SuggestedFile: FC<ISuggestedFileProps> = (props) => {
  const { file } = props;

  const classes = useStyles();
  const colorCode = ColorUtils.getColorCode(file.extension);

  return (
    <Box
      className={classes.singleFileWrapper}
      style={{
        background: colorCode.backgroundGradient
      }}
    >
      <Box
        display={'flex'}
        flexDirection={'column'}
        width={'80%'}
        className={'truncate'}
      >
        <Box mb={1} display={'flex'} alignItems={'center'} justifyContent={'center'}>

          <FileIcon
            extension={file.extension}
            background={'transparent'}
            colorScheme={{ primary: colorCode.iconColor }}
            iconStyle={IconStyle.normal}
            size={45}
          />
        </Box>
          <Typography
            className={clsx('truncate', classes.fileName)}
          >{`${file.name}.${file.extension}`}</Typography>
      </Box>
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    singleFileWrapper: {
      width: '160px',
      height: '150px',
      borderRadius: '15px',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      background: theme.palette.workspaceGradients.lightBrown,
      transition: 'box-shadow .24s ease-in-out',
      cursor: 'pointer',
      '&:hover': {
        boxShadow: theme.palette.boxShadows.darker
      },
      wordWrap: 'break-word'
    },
    fileName: {
      fontWeight: 500,
      fontSize: theme.typography.pxToRem(12)
    }
  };
});

export default SuggestedFile;
