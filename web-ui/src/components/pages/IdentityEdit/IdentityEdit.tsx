import {
  Avatar,
  Box,
  IconButton,
  TextField,
  Typography
} from '@material-ui/core';
import AddRoundedIcon from '@material-ui/icons/AddRounded';
import ArrowBackRoundedIcon from '@material-ui/icons/ArrowBackRounded';
import CameraAltRoundedIcon from '@material-ui/icons/CameraAltRounded';
import SaveRoundedIcon from '@material-ui/icons/SaveRounded';
import { useFormik } from 'formik';
import { createRef, FC, useCallback, useEffect, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import { useHistory, useParams } from 'react-router-dom';
import IdentitiesService from '../../../services/identities/identitiesService';
import {
  INewIdentityRequest,
  IUpdateIdentityRequest
} from '../../../services/identities/identitiesService.types';
import newContactValidationSchema from '../../../shared/schemas/contactSchemas';
import CommonUtils from '../../../shared/utils/CommonUtils';
import theme from '../../../theme/theme';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import Link from '../../atoms/Link/Link';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import useSnackbar from '../../molecules/Snackbar/useSnackbar.hook';
import KeyManager from '../../organisms/KeyManager/KeyManager';
import {
  EKeyInputType,
  EKeyType
} from '../../organisms/KeyManager/keyManager.types';
import { IKeyPair } from '../ContactEdit/contactEdit.types';

import {
  EIdentityEditType,
  IIdentityEditParams,
  IIdentityEditProps
} from './identityEdit.types';

const IdentityEdit: FC<IIdentityEditProps> = (props) => {
  const { type } = props;

  const pageTitle =
    type === EIdentityEditType.NEW ? 'New Identity' : 'Edit Identity';

  const [pictureSrc, setPictureSrc] = useState<string>('');

  const [initialValues, setInitialValues] = useState<{
    name: string;
    picture: string;
    keyPair: IKeyPair;
  }>({
    name: '',
    picture: '',
    keyPair: null
  });

  const handleNewIdentity = async (request: INewIdentityRequest) => {
    await IdentitiesService.createIdentity(request);
  };

  const handleUpdateIdentity = async (request: IUpdateIdentityRequest) => {
    await IdentitiesService.updateIdentity(request);
  };

  const history = useHistory();

  const formik = useFormik({
    initialValues,
    enableReinitialize: true,
    validationSchema: newContactValidationSchema,
    onSubmit: (values, { resetForm }) => {
      setErrorMessage('');

      if (values.keyPair == null) {
        setErrorMessage('A key pair is required');
        return;
      }

      if (type == EIdentityEditType.NEW) {
        handleNewIdentity({
          name: values.name,
          picture: values.picture,
          privateKey: values.keyPair.privateKey
        })
          .then((response) => {
            openSnackbar('Identity successfully created', 'success');
            history.push('/identities');
          })
          .catch((err) => {
            resetForm();

            openSnackbar('Unable to create identity', 'error');
          });
      } else {
        handleUpdateIdentity({
          identityId,
          name: values.name,
          picture: values.picture,
          privateKey: values.keyPair.privateKey
        })
          .then((response) => {
            openSnackbar('Identity successfully updated', 'success');
            history.push('/identities');
          })
          .catch((err) => {
            resetForm();

            openSnackbar('Unable to create identity', 'error');
          });
      }
    }
  });

  const dropzoneRef: any = createRef();
  const openDialog = () => {
    // Note that the ref is set async,
    // so it might be null at some point
    if (dropzoneRef.current) {
      dropzoneRef.current.open();
    }
  };

  const onDrop = useCallback(async (acceptedFiles: Array<File>) => {
    if (acceptedFiles.length > 0) {
      let base64 = await CommonUtils.fileToBase64(acceptedFiles[0]);

      await formik.setFieldValue('picture', base64);
      setPictureSrc(base64);
    }
  }, []);

  const {
    getRootProps,
    getInputProps,
    isDragActive,
    isDragAccept,
    isDragReject
  } = useDropzone({ onDrop, accept: '.jpg, .jpeg, .png', maxFiles: 1 });

  const [errorMessage, setErrorMessage] = useState<string>();
  const [addedKey, setAddedKey] = useState<IKeyPair>(null);

  const { identityId } = useParams() as IIdentityEditParams;
  const { openSnackbar } = useSnackbar();

  useEffect(() => {
    if (type === EIdentityEditType.EDIT) {
      // Fetch contact information based on the passed in ID
      const fetchIdentity = async () => {
        return await IdentitiesService.getIdentity(identityId);
      };

      fetchIdentity()
        .then((response) => {
          let keyPair = {
            keyID: response.publicKeyID,
            publicKey: '',
            privateKey: ''
          };

          setInitialValues({
            name: response.name,
            picture: response.picture,
            keyPair
          });
          setAddedKey(keyPair);
          setPictureSrc(response.picture);
        })
        .catch((err) => {
          openSnackbar('Unable to fetch identity', 'error');

          setAddedKey(null);
          setInitialValues({
            name: '',
            picture: '',
            keyPair: null
          });
        });
    }
  }, []);

  return (
    <Box display={'flex'} flexDirection={'column'}>
      <Box display={'flex'} alignItems={'center'}>
        {
          // TODO add unsaved changes modal
        }
        <Link to={'/identities'}>
          <IconButton
            classes={{
              root: 'iconButtonRoot'
            }}
          >
            <ArrowBackRoundedIcon
              style={{
                fill: 'black'
              }}
            />
          </IconButton>
        </Link>
        <Box ml={2}>
          <PageTitle title={pageTitle} />
        </Box>
      </Box>
      <Box width={'100%'} display={'flex'} flexDirection={'column'}>
        <form autoComplete={'off'} onSubmit={formik.handleSubmit}>
          <Box display={'flex'} maxWidth={'100%'} mt={2}>
            <Box display={'flex'} flexDirection={'column'} maxWidth={'30%'}>
              <Box mb={2}>
                <FormTitle title={'Basic info'} />
              </Box>
              <TextField
                id={'name'}
                label={'Name'}
                variant={'outlined'}
                value={formik.values.name}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.touched.name && Boolean(formik.errors.name)}
                helperText={formik.touched.name && formik.errors.name}
              />
            </Box>
            <Box
              display={'flex'}
              flexDirection={'column'}
              width={'auto'}
              alignItems={'center'}
              justifyContent={'center'}
              ml={6}
            >
              <Avatar
                src={pictureSrc}
                style={{
                  width: '70px',
                  height: '70px',
                  marginBottom: '16px'
                }}
              />
              <Box component={'div'} {...getRootProps()}>
                <input {...getInputProps()} />
                <ActionButton
                  text={'Select image'}
                  startIcon={<CameraAltRoundedIcon />}
                  shouldSubmit={false}
                  onClick={openDialog}
                />
              </Box>
            </Box>
          </Box>
          <Box display={'flex'} flexDirection={'column'} mt={2}>
            <Box mb={2}>
              <FormTitle title={'Key pair'} />
            </Box>
            <KeyManager
              addedKey={addedKey}
              expectedType={EKeyType.PRIVATE}
              keyListTitle={
                type == EIdentityEditType.EDIT
                  ? 'Existing pair'
                  : 'Created pair'
              }
              setAddedKey={setAddedKey}
              visibleTypes={[
                EKeyInputType.IMPORT,
                EKeyInputType.ENTER,
                EKeyInputType.GENERATE
              ]}
              formik={formik}
            />
          </Box>
          <Box display={errorMessage ? 'flex' : 'none'} mt={2}>
            <Typography
              variant={'body1'}
              style={{
                marginTop: '1rem',
                fontFamily: 'Montserrat',
                color: theme.palette.error.main
              }}
            >
              {errorMessage}
            </Typography>
          </Box>
          <Box mt={5} display={'flex'}>
            <ActionButton
              text={type == EIdentityEditType.EDIT ? 'Save' : 'Create'}
              square
              startIcon={
                type == EIdentityEditType.EDIT ? (
                  <SaveRoundedIcon />
                ) : (
                  <AddRoundedIcon />
                )
              }
            />
          </Box>
        </form>
      </Box>
    </Box>
  );
};

export default IdentityEdit;
