import {
  Box,
  FormControlLabel,
  Radio,
  RadioGroup,
  TextField,
  Typography
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import BackupRoundedIcon from '@material-ui/icons/BackupRounded';
import { FC, Fragment, useMemo, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import theme from '../../../theme/theme';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import KeyList from '../../atoms/KeyList/KeyList';
import { EKeyInputType, IKeyManagerProps } from './keyManager.types';

const KeyManager: FC<IKeyManagerProps> = (props) => {
  const [activeTab, setActiveTab] = useState<EKeyInputType>(
    EKeyInputType.IMPORT
  );
  const classes = useStyles();

  const dropzoneStyle = {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    padding: '20px',
    height: '100%',
    borderWidth: 2,
    borderRadius: 2,
    borderColor: '#eeeeee',
    borderStyle: 'dashed',
    backgroundColor: '#fafafa',
    color: '#bdbdbd',
    outline: 'none',
    transition: 'border .24s ease-in-out'
  };

  const {
    getRootProps,
    getInputProps,
    isDragActive,
    isDragAccept,
    isDragReject
  } = useDropzone({ accept: 'image/*' });

  const style = useMemo(
    () => ({
      ...dropzoneStyle
    }),
    [isDragActive, isDragReject, isDragAccept]
  );

  const handleTypeSelect = (type: EKeyInputType) => {
    setActiveTab(type);
  };

  const { visibleTypes } = props;

  const [addedKeys, setAddedKeys] = useState<{ keys: string[] }>({
    keys: []
  });

  const handleKeyRemove = (key: string) => {
    const index = addedKeys.keys.indexOf(key);
    if (index > -1) {
      let newArr = addedKeys.keys;
      newArr.splice(index, 1);
      setAddedKeys({ keys: newArr });
    }
  };

  return (
    <Box display={'flex'} flexDirection={'column'} width={'100%'}>
      {
        // Radio selects
      }
      <Box display={'flex'} width={'100%'}>
        <RadioGroup name="keyInputType" value={activeTab} row>
          {visibleTypes.map((type, index) => {
            return (
              <FormControlLabel
                value={type}
                style={{
                  marginLeft: 0
                }}
                control={<Radio color={'primary'} />}
                label={type}
                onChange={() => {
                  handleTypeSelect(type);
                }}
              />
            );
          })}
        </RadioGroup>
      </Box>
      {
        // Dropzone and added keys
      }

      <Box mt={2} display={'flex'}>
        {activeTab === EKeyInputType.IMPORT && (
          <Fragment>
            <Box className={classes.dropzoneWrapper}>
              <Box component={'div'} {...getRootProps(style)}>
                <input {...getInputProps()} />
                <Box
                  display={'flex'}
                  flexDirection={'column'}
                  width={'100%'}
                  height={'100%'}
                  alignItems={'center'}
                  justifyContent={'center'}
                >
                  <BackupRoundedIcon />
                  <Typography>
                    <span
                      style={{
                        fontWeight: 600
                      }}
                    >
                      Click to upload
                    </span>{' '}
                    or drag and drop
                  </Typography>
                </Box>
              </Box>
            </Box>

            <Box
              display={'flex'}
              ml={4}
              flexDirection={'column'}
              maxHeight={'200px'}
            >
              <FormTitle title={'Added keys'} />
              <KeyList
                addedKeys={addedKeys.keys}
                handleKeyRemove={handleKeyRemove}
              />
            </Box>
          </Fragment>
        )}

        {activeTab === EKeyInputType.ENTER && (
          <Fragment>
            <Box className={classes.textInputWrapper}>
              {
                // TODO add on input change listener
              }
              <TextField
                variant={'outlined'}
                rows={10}
                style={{
                  width: '100%',
                  fontSize: '0.875rem'
                }}
                size={'small'}
                multiline
                placeholder={
                  '-----BEGIN PGP PUBLIC KEY BLOCK-----\n' +
                  'Version: GnuPG v1.2.1 (GNU/Linux)\n\n' +
                  '...\n\n' +
                  '-----END PGP PUBLIC KEY BLOCK-----'
                }
              />
            </Box>

            <Box
              display={'flex'}
              ml={4}
              flexDirection={'column'}
              maxHeight={'200px'}
            >
              <FormTitle title={'Added keys'} />
              <KeyList
                addedKeys={addedKeys.keys}
                handleKeyRemove={handleKeyRemove}
              />
            </Box>
          </Fragment>
        )}
      </Box>
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    dropzoneWrapper: {
      backgroundColor: theme.palette.custom.darkGray,
      borderStyle: 'dashed',
      border: '2px solid #A7A7A7',
      borderRadius: '15px',
      height: '200px',
      cursor: 'pointer',
      width: '50%'
    },
    textInputWrapper: {
      borderRadius: '15px',
      height: '200px',
      width: '50%',
      fontSize: '0.875rem'
    }
  };
});
export default KeyManager;
