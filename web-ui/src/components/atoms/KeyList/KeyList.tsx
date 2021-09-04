import { Box, IconButton, Typography } from '@material-ui/core';
import CloseRoundedIcon from '@material-ui/icons/CloseRounded';
import VpnKeyRoundedIcon from '@material-ui/icons/VpnKeyRounded';
import { FC } from 'react';
import { IKeyListProps } from './keyList.types';

const KeyList: FC<IKeyListProps> = (props) => {
  const { addedKeys, handleKeyRemove } = props;

  if (addedKeys.length == 0) {
    return <Typography>No keys added</Typography>;
  } else {
    return (
      <Box
        display={'flex'}
        flexDirection={'column'}
        flexWrap={'wrap'}
        height={'100%'}
        marginLeft={'-15px'}
      >
        {addedKeys.map((addedKey) => {
          return (
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
                    {addedKey}
                  </Typography>
                </Box>
              </Box>

              <Box ml={'8px'}>
                <IconButton onClick={() => handleKeyRemove(addedKey)}>
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
          );
        })}
      </Box>
    );
  }
};

export default KeyList;
