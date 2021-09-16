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
import { Formik } from 'formik';
import { FC, Fragment, useCallback, useMemo, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import CryptoService from '../../../services/crypto/cryptoService';
import { IValidatePublicKeyResponse } from '../../../services/crypto/cryptoService.types';
import generateKeyValidationSchema from '../../../shared/schemas/identitySchemas';
import CommonUtils from '../../../shared/utils/CommonUtils';
import theme from '../../../theme/theme';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import KeyList from '../../atoms/KeyList/KeyList';
import IdentityOverwriteModal from '../../molecules/IdentityOverwriteModal/IdentityOverwriteModal';
import useSnackbar from '../../molecules/Snackbar/useSnackbar.hook';
import { IKeyPair } from '../../pages/ContactEdit/contactEdit.types';
import {
  EKeyGenerateType,
  EKeyInputType,
  EKeyType,
  IKeyManagerProps
} from './keyManager.types';

const KeyManager: FC<IKeyManagerProps> = (props) => {
  const [activeTab, setActiveTab] = useState<EKeyInputType>(
    EKeyInputType.IMPORT
  );

  const [activeKeyGeneration, setActiveKeyGeneration] =
    useState<EKeyGenerateType>(EKeyGenerateType.RSA_2048);

  const [generatedKey, setGeneratedKey] = useState<string>('');

  const classes = useStyles();

  const {
    addedKey,
    setAddedKey,
    keyListTitle = 'Added key',
    formik,
    expectedType
  } = props;

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

  const { openSnackbar } = useSnackbar();

  const parseKey = async (key: string) => {
    try {
      if (expectedType == EKeyType.PUBLIC) {
        let response: IValidatePublicKeyResponse =
          await CryptoService.validatePublicKey(key);

        await updateKey({
          keyID: response.publicKeyID,
          publicKey: key,
          privateKey: ''
        });
      } else {
        // TODO
      }

      openSnackbar('Key successfully added', 'success');
    } catch (error) {
      openSnackbar('Unable to add key', 'error');
    }
  };

  const onDrop = useCallback((acceptedFiles) => {
    const reader = new FileReader();
    reader.onload = async (e) => {
      const text = e.target.result;
      const preparedKey: string = CommonUtils.removeLineBreaks(text as string);

      await parseKey(preparedKey);
    };

    reader.readAsText(acceptedFiles[0]);
  }, []);

  const {
    getRootProps,
    getInputProps,
    isDragActive,
    isDragAccept,
    isDragReject
  } = useDropzone({ onDrop, accept: '.asc, .gpg, .pgp', maxFiles: 1 });

  const style = useMemo(
    () => ({
      ...dropzoneStyle
    }),
    [isDragActive, isDragReject, isDragAccept]
  );

  const [keyInputValue, setKeyInputValue] = useState<string>('');

  const handleTypeSelect = (type: EKeyInputType) => {
    setActiveTab(type);
  };

  const handleKeyTypeSelect = (type: EKeyGenerateType) => {
    setActiveKeyGeneration(type);
  };

  const { visibleTypes } = props;

  const handleKeyRemove = () => {
    setAddedKey(null);
    setBufferKeyPair(null);
    setKeyInputValue('');
  };

  const validKeyTypes = [EKeyGenerateType.RSA_2048, EKeyGenerateType.RSA_4096];

  const dummyPK =
    '-----BEGIN RSA PRIVATE KEY-----\n' +
    'MIICXAIBAAKBgQCqGKukO1De7zhZj6+H0qtjTkVxwTCpvKe4eCZ0FPqri0cb2JZfXJ/DgYSF6vUp' +
    'wmJG8wVQZKjeGcjDOL5UlsuusFncCzWBQ7RKNUSesmQRMSGkVb1/3j+skZ6UtW+5u09lHNsj6tQ5' +
    '1s1SPrCBkedbNf0Tp0GbMJDyR4e9T04ZZwIDAQABAoGAFijko56+qGyN8M0RVyaRAXz++xTqHBLh' +
    '3tx4VgMtrQ+WEgCjhoTwo23KMBAuJGSYnRmoBZM3lMfTKevIkAidPExvYCdm5dYq3XToLkkLv5L2' +
    'pIIVOFMDG+KESnAFV7l2c+cnzRMW0+b6f8mR1CJzZuxVLL6Q02fvLi55/mbSYxECQQDeAw6fiIQX' +
    'GukBI4eMZZt4nscy2o12KyYner3VpoeE+Np2q+Z3pvAMd/aNzQ/W9WaI+NRfcxUJrmfPwIGm63il' +
    'AkEAxCL5HQb2bQr4ByorcMWm/hEP2MZzROV73yF41hPsRC9m66KrheO9HPTJuo3/9s5p+sqGxOlF' +
    'L0NDt4SkosjgGwJAFklyR1uZ/wPJjj611cdBcztlPdqoxssQGnh85BzCj/u3WqBpE2vjvyyvyI5k' +
    'X6zk7S0ljKtt2jny2+00VsBerQJBAJGC1Mg5Oydo5NwD6BiROrPxGo2bpTbu/fhrT8ebHkTz2epl' +
    'U9VQQSQzY1oZMVX8i1m5WUTLPz2yLJIBQVdXqhMCQBGoiuSoSjafUhV7i1cEGpb88h5NBYZzWXGZ' +
    '37sJ5QsW+sJyoNde3xH8vdXhzU7eT82D6X/scw9RZz+/6rCJ4p0=\n' +
    '-----END RSA PRIVATE KEY-----';

  const handleGenerateKey = async () => {
    let keyPair = {
      keyID: '4AEE18F83AFDEB23',
      publicKey: '4AEE18F83AFDEB23',
      privateKey: '4AEE18F83AFDEB23'
    };

    setBufferKeyPair(keyPair);

    if (addedKey != null) {
      setConfirmOpen(true);
    } else {
      await updateKey(keyPair);
    }
  };

  // Overwrite modal
  const handleConfirm = async (confirmed: boolean) => {
    console.log(confirmed);
    if (confirmed) {
      console.log('Overwriting');
      await updateKey(bufferKeyPair);
    }

    setBufferKeyPair(null);
    setConfirmOpen(false);
  };

  const updateKey = async (keyPair: IKeyPair) => {
    setGeneratedKey(dummyPK); // TODO swap out
    setAddedKey(keyPair);

    await formik.setFieldValue('keyPair', keyPair);
    console.log(formik.values.keyPair);
  };

  const [confirmOpen, setConfirmOpen] = useState<boolean>(false);

  // Used for state management in overwriting
  const [bufferKeyPair, setBufferKeyPair] = useState<IKeyPair>(null);
  const handleKeyInput = async (event: any) => {
    setKeyInputValue(event.target.value);

    // Trigger validation and addition
    await parseKey(event.target.value);
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

      <Box mt={4} display={'flex'}>
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
              <FormTitle title={keyListTitle} />
              <KeyList addedKey={addedKey} handleKeyRemove={handleKeyRemove} />
            </Box>
          </Fragment>
        )}

        {activeTab === EKeyInputType.ENTER && (
          <Fragment>
            <Box className={classes.textInputWrapper}>
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
                onChange={handleKeyInput}
                onKeyUp={handleKeyInput}
                onPaste={handleKeyInput}
                onInput={handleKeyInput}
                value={keyInputValue}
              />
            </Box>

            <Box
              display={'flex'}
              ml={4}
              flexDirection={'column'}
              maxHeight={'200px'}
            >
              <FormTitle title={keyListTitle} />
              <KeyList addedKey={addedKey} handleKeyRemove={handleKeyRemove} />
            </Box>
          </Fragment>
        )}
        {activeTab === EKeyInputType.GENERATE && (
          <Fragment>
            <Box display={'flex'} flexDirection={'column'} width={'100%'}>
              <Formik
                initialValues={{
                  generateName: '',
                  generateEmail: ''
                }}
                enableReinitialize={true}
                validationSchema={generateKeyValidationSchema}
                onSubmit={(values, { resetForm }) => {
                  handleGenerateKey();
                }}
              >
                {(keyGenerationFormik) => (
                  <form autoComplete={'off'}>
                    <Box width={'100%'} display={'flex'} alignItems={'center'}>
                      <Box minHeight={'80px'}>
                        <TextField
                          id={'generateName'}
                          label={'Name'}
                          variant={'outlined'}
                          value={keyGenerationFormik.values.generateName}
                          onChange={keyGenerationFormik.handleChange}
                          onBlur={keyGenerationFormik.handleBlur}
                          error={
                            keyGenerationFormik.touched.generateName &&
                            Boolean(keyGenerationFormik.errors.generateName)
                          }
                          helperText={
                            keyGenerationFormik.touched.generateName &&
                            keyGenerationFormik.errors.generateName
                          }
                        />
                      </Box>

                      <Box ml={4} minHeight={'80px'}>
                        <TextField
                          id={'generateEmail'}
                          label={'Email'}
                          variant={'outlined'}
                          value={keyGenerationFormik.values.generateEmail}
                          onChange={keyGenerationFormik.handleChange}
                          onBlur={keyGenerationFormik.handleBlur}
                          error={
                            keyGenerationFormik.touched.generateEmail &&
                            Boolean(keyGenerationFormik.errors.generateEmail)
                          }
                          helperText={
                            keyGenerationFormik.touched.generateEmail &&
                            keyGenerationFormik.errors.generateEmail
                          }
                        />
                      </Box>
                    </Box>

                    <Box display={'flex'} width={'100%'} mt={2}>
                      <Box display={'flex'} width={'50%'}>
                        <RadioGroup
                          name="keyGenerationType"
                          value={activeKeyGeneration}
                          row
                        >
                          {validKeyTypes.map((type) => {
                            return (
                              <FormControlLabel
                                key={'keyType-' + type}
                                value={type}
                                style={{
                                  marginLeft: 0
                                }}
                                control={<Radio color={'primary'} />}
                                label={type}
                                onChange={() => {
                                  handleKeyTypeSelect(type);
                                }}
                              />
                            );
                          })}
                        </RadioGroup>
                        <Box ml={'auto'}>
                          <ActionButton
                            square={true}
                            text={'Generate'}
                            shouldSubmit={false}
                            onClick={() => keyGenerationFormik.submitForm()}
                          />
                        </Box>
                      </Box>
                      <Box
                        display={'flex'}
                        ml={4}
                        flexDirection={'column'}
                        maxHeight={'200px'}
                      >
                        <FormTitle title={keyListTitle} />
                        <KeyList
                          addedKey={addedKey}
                          handleKeyRemove={() => {
                            handleKeyRemove();
                            setGeneratedKey('');
                            formik.values.keyPair = null;
                          }}
                        />
                      </Box>
                    </Box>

                    <Box className={classes.textInputWrapper} mt={4}>
                      {
                        // TODO add on input change listener
                      }
                      <TextField
                        variant={'outlined'}
                        rows={10}
                        value={generatedKey}
                        contentEditable={false}
                        disabled={false}
                        style={{
                          width: '100%',
                          fontSize: '0.875rem'
                        }}
                        size={'small'}
                        multiline
                        placeholder={
                          '-----BEGIN PGP PRIVATE KEY BLOCK-----\n' +
                          'Version: GnuPG v1.2.1 (GNU/Linux)\n\n' +
                          '...\n\n' +
                          '-----END PGP PRIVATE KEY BLOCK-----'
                        }
                      />
                    </Box>
                  </form>
                )}
              </Formik>
            </Box>
          </Fragment>
        )}
      </Box>
      <IdentityOverwriteModal
        publicKeyID={addedKey ? addedKey.keyID : ''}
        open={confirmOpen}
        handleConfirm={handleConfirm}
      />
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
    inputRoot: {
      fontSize: '0.5rem'
    },
    textInputWrapper: {
      borderRadius: '15px',
      height: '200px',
      width: '50%',
      fontSize: '0.875rem !important'
    }
  };
});
export default KeyManager;
