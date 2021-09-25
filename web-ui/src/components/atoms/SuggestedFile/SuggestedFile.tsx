import { Box, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import clsx from 'clsx';
import { saveAs } from 'file-saver';
import { FC } from 'react';
import FileIcon, { IconStyle } from 'react-fileicons';
import WorkspacesService from '../../../services/workspaces/workspacesService';
import ColorUtils from '../../../shared/utils/ColorUtils';
import CommonUtils from '../../../shared/utils/CommonUtils';
import theme from '../../../theme/theme';
import useSnackbar from '../../molecules/Snackbar/useSnackbar.hook';
import { ISuggestedFileProps } from './suggestedFile.types';

const SuggestedFile: FC<ISuggestedFileProps> = (props) => {
  const { file, workspaceMnemonic } = props;

  const classes = useStyles();
  const colorCode = ColorUtils.getColorCode(file.extension);

  const { openSnackbar } = useSnackbar();
  const renderItemName = (name: string, extension: string) => {
    return `${name}${extension}`;
  };

  const handleDownload = (
    checksum: string,
    name: string,
    extension: string
  ) => {
    openSnackbar('Started file download...', 'success');

    const downloadFile = async () => {
      return await WorkspacesService.downloadFile({
        workspaceMnemonic: CommonUtils.unformatMnemonic(workspaceMnemonic),
        fileChecksum: checksum
      });
    };

    downloadFile()
      .then((response) => {
        saveAs(response, renderItemName(name, extension));

        openSnackbar('File successfully downloaded', 'success');
      })
      .catch((err) => {
        openSnackbar('Unable to download file', 'error');
      });
  };

  return (
    <Box
      className={classes.singleFileWrapper}
      style={{
        background: colorCode.backgroundGradient
      }}
      ml={4}
      onClick={() => {
        handleDownload(file.checksum, file.name, file.extension);
      }}
    >
      <Box
        display={'flex'}
        flexDirection={'column'}
        width={'80%'}
        className={'truncate'}
      >
        <Box
          mb={1}
          display={'flex'}
          alignItems={'center'}
          justifyContent={'center'}
        >
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
        >{`${file.name}${file.extension}`}</Typography>
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
      textAlign: 'center',
      fontWeight: 500,
      fontSize: theme.typography.pxToRem(12)
    }
  };
});

export default SuggestedFile;
