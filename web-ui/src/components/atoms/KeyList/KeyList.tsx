import { Box, IconButton, Typography } from '@material-ui/core';
import CloseRoundedIcon from '@material-ui/icons/CloseRounded';
import VpnKeyRoundedIcon from '@material-ui/icons/VpnKeyRounded';
import { FC } from 'react';
import { IKeyListProps } from './keyList.types';

const KeyList: FC<IKeyListProps> = (props) => {
  const { addedKey, handleKeyRemove } = props;

  if (addedKey == null) {
    return <Typography>No key pair added</Typography>;
  } else {
    return (
      <Box
        display={'flex'}
        flexDirection={'column'}
        flexWrap={'wrap'}
        height={'100%'}
        marginLeft={'-15px'}
      >
        <Box
          display={'flex'}
          alignItems={'center'}
          marginLeft={'15px'}
          justifyContent={'space-between'}
        >
          <Box display={'flex'} alignItems={'center'}>
            <VpnKeyRoundedIcon />
            <Box ml={1}>
              <Typography style={{ fontSize: '0.875rem' }}>
                {addedKey.keyID}
              </Typography>
            </Box>
          </Box>

          <Box ml={'8px'}>
            <IconButton
              classes={{
                root: 'iconButtonRoot'
              }}
              onClick={() => handleKeyRemove(addedKey)}
            >
              <CloseRoundedIcon
                style={{
                  fill: 'black',
                  width: '20px',
                  height: 'auto'
                }}
              />
            </IconButton>
          </Box>
        </Box>
      </Box>
    );
  }
};

export default KeyList;
