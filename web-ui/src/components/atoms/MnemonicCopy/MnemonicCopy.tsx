import { Box, IconButton, Tooltip, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { FC, useState } from 'react';
import { ReactComponent as Clipboard } from '../../../shared/assets/icons/content_paste_black_24dp.svg';
import theme from '../../../theme/theme';
import { IMnemonicCopyProps } from './mnemonicCopy.types';

const MnemonicCopy: FC<IMnemonicCopyProps> = (props) => {
  const { mnemonic } = props;

  const classes = useStyles();

  const [clipboardText, setClipboardText] =
    useState<string>('Copy to clipboard');

  const handleCopy = () => {
    navigator.clipboard.writeText(mnemonic);

    setClipboardText('Copied to clipboard!');
  };

  return (
    <Box className={classes.mnemonicWrapper}>
      <Typography className={classes.mnemonicText}>{mnemonic}</Typography>
      <Box ml={4}>
        <Tooltip
          title={clipboardText}
          arrow
          onClose={() => {
            setClipboardText('Copy to clipboard');
          }}
        >
          <IconButton
            classes={{
              root: 'iconButtonRoot'
            }}
            onClick={handleCopy}
          >
            <Clipboard />
          </IconButton>
        </Tooltip>
      </Box>
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    mnemonicText: {
      fontWeight: 600,
      fontSize: theme.typography.pxToRem(22)
    },
    mnemonicWrapper: {
      borderRadius: '15px',
      backgroundColor: '#F0F0F0',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '20px 25px',
      height: '80px',
      width: '50%'
    }
  };
});

export default MnemonicCopy;
