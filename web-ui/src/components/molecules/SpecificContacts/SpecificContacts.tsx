import {
  Box,
  Checkbox,
  Dialog,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  IconButton,
  TextField,
  Typography
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import AddRoundedIcon from '@material-ui/icons/AddRounded';
import CloseRoundedIcon from '@material-ui/icons/CloseRounded';
import { FC, useEffect, useState } from 'react';
import { ContactResponse } from '../../../context/newWorkspaceContext.types';
import theme from '../../../theme/theme';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import { ISpecificContactsProps } from './specificContacts.types';

const SpecificContacts: FC<ISpecificContactsProps> = (props) => {
  const {
    listTitle = 'Permitted contacts',
    disabled = false,
    contactIDs,
    setContactIDs
  } = props;

  const classes = useStyles();
  const [dialogOpen, setDialogOpen] = useState(false);

  const { openSnackbar } = useSnackbar();

  const [loading, setLoading] = useState<boolean>(true);

  const [allContacts, setAllContacts] = useState<
    { list: ContactResponse[] } | undefined
  >();

  const [filteredList, setFilteredList] = useState<
    { list: ContactResponse[] } | undefined
  >({ list: [] });

  const [selectedContacts, setSelectedContacts] = useState<{
    list: ContactResponse[];
  }>({
    list: contactIDs.contacts
  });

  useEffect(() => {
    if (dialogOpen) {
      const fetchAllContacts = async () => {
        return [
          {
            id: '123',
            name: 'Milos Zivkovic',
            publicKeyID: '123'
          },
          {
            id: '1234',
            name: 'Milos Zivkovic',
            publicKeyID: '1234'
          },
          {
            id: '12345',
            name: 'Milos Zivkovic',
            publicKeyID: '12345'
          },
          {
            id: '123456',
            name: 'Milos Zivkovic',
            publicKeyID: '123456'
          },
          {
            id: '1234567',
            name: 'Milos Zivkovic',
            publicKeyID: '123456'
          },
          {
            id: '12345678',
            name: 'Milos Zivkovic',
            publicKeyID: '123456'
          },
          {
            id: '12345689',
            name: 'Milos Zivkovic',
            publicKeyID: '123456'
          },
          {
            id: '1234568910',
            name: 'Milos Zivkovic',
            publicKeyID: '123456'
          },
          {
            id: '1234568911',
            name: 'Milos Zivkovic',
            publicKeyID: '123456'
          }
        ];
      };

      fetchAllContacts()
        .then((response) => {
          setAllContacts({ list: response });

          setFilteredList({ list: response });
        })
        .catch((err) => {
          openSnackbar('Unable to fetch contact list', 'error');
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, [dialogOpen]);

  const isContactChecked = (id: string) => {
    return selectedContacts.list.some((itemId) => itemId.id === id);
  };

  const findContact = (id: string) => {
    for (let i = 0; i < selectedContacts.list.length; i++) {
      if (selectedContacts.list[i].id === id) return i;
    }

    return -1;
  };

  const handleContactToggle = async (
    id: string,
    name: string,
    publicKeyID: string
  ) => {
    const index = findContact(id);

    if (index > -1) {
      let newArr = selectedContacts.list;
      newArr.splice(index, 1);

      setSelectedContacts({ list: newArr });
    } else {
      let newArr = selectedContacts.list;
      newArr.push({ id, name, publicKeyID });

      setSelectedContacts({ list: newArr });
    }
  };

  const manualFilter = (fn: any, array: any) => {
    const f = [];

    for (let i = 0; i < array.length; i++) {
      if (fn(array[i])) {
        f.push(array[i]);
      }
    }
    return f;
  };

  const handleSearchInputChange = (event: any) => {
    const value = event.target.value;

    let tempFilteredContacts = { list: allContacts.list };
    // Filter if there is a value in search input
    if (value) {
      tempFilteredContacts.list = manualFilter(
        (contact: any) =>
          contact.name.toLowerCase().includes(value.toLowerCase()) ||
          contact.publicKeyID.toLowerCase().includes(value.toLowerCase()),
        allContacts.list
      );
    }

    setFilteredList(tempFilteredContacts);
  };

  const handleDialogToggle = (open: boolean) => {
    setDialogOpen(open);

    let contactIDs: ContactResponse[] = [];
    for (let i = 0; i < selectedContacts.list.length; i++) {
      contactIDs.push(selectedContacts.list[i]);
    }

    setContactIDs({ contacts: contactIDs });
  };

  return (
    <Box display={'flex'} flexDirection={'column'}>
      <Box display={'flex'} alignItems={'center'} mb={2}>
        <FormTitle title={listTitle} />
        <Box ml={4}>
          <ActionButton
            text={'Select contacts'}
            disabled={disabled}
            shouldSubmit={false}
            startIcon={<AddRoundedIcon />}
            onClick={() => {
              handleDialogToggle(true);
            }}
          />
        </Box>
      </Box>
      <Box>
        {selectedContacts.list.length < 1 ? (
          <Typography className={classes.noContacts}>
            No contacts added
          </Typography>
        ) : (
          <Box
            mb={-1}
            ml={-1}
            display={'flex'}
            flexDirection={'column'}
            flexWrap={'wrap'}
            maxHeight={'200px'}
            minHeight={'200px'}
          >
            {selectedContacts.list.map(
              (selectedContact: ContactResponse, index) => {
                return (
                  <Box mb={1} ml={1}>
                    <Typography>{`${index + 1}. ${selectedContact.name} (${
                      selectedContact.publicKeyID
                    })`}</Typography>
                  </Box>
                );
              }
            )}
          </Box>
        )}
      </Box>

      <Box display={props.errorMessage ? 'flex' : 'none'}>
        <Typography
          variant={'body1'}
          style={{
            marginTop: '1rem',
            fontFamily: 'Montserrat',
            color: theme.palette.error.main
          }}
        >
          {props.errorMessage}
        </Typography>
      </Box>
      <Dialog
        disableBackdropClick
        disableEscapeKeyDown
        open={dialogOpen}
        onClose={() => handleDialogToggle(false)}
        classes={{ paper: classes.dialogPaper }}
      >
        <DialogTitle disableTypography>
          <Box display={'flex'} alignItems={'center'}>
            <Typography className={classes.dialogTitle}>
              Choose contacts
            </Typography>
            <Box ml={'auto'}>
              <IconButton
                classes={{
                  root: 'iconButtonRoot'
                }}
                onClick={() => handleDialogToggle(false)}
              >
                <CloseRoundedIcon />
              </IconButton>
            </Box>
          </Box>
        </DialogTitle>
        <DialogContent>
          <Box mb={2}>
            <TextField
              style={{
                width: '100%'
              }}
              placeholder={'Search'}
              onChange={handleSearchInputChange}
            />
          </Box>

          {filteredList &&
            filteredList.list.map((contact) => {
              return (
                <FormControlLabel
                  key={`contact-${contact.id}`}
                  value={contact.id}
                  control={<Checkbox color={'primary'} />}
                  label={
                    <Typography style={{ textAlign: 'left' }}>
                      {`${contact.name} (${contact.publicKeyID})`}
                    </Typography>
                  }
                  labelPlacement="end"
                  checked={isContactChecked(contact.id)}
                  onChange={() =>
                    handleContactToggle(
                      contact.id,
                      contact.name,
                      contact.publicKeyID
                    )
                  }
                />
              );
            })}
          {loading && <LoadingIndicator size={54} />}
          {!loading && filteredList && filteredList.list.length < 1 && (
            <Box
              display={'flex'}
              width={'100%'}
              justifyContent={'center'}
              mt={1}
            >
              <Typography>No contacts found</Typography>
            </Box>
          )}
        </DialogContent>
      </Dialog>
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    noContacts: {
      fontSize: theme.typography.pxToRem(16),
      color: theme.palette.custom.transparentBlack
    },
    dialogPaper: {
      height: '100%',
      maxHeight: '100%',
      width: '100%',
      margin: 0,
      [theme.breakpoints.up('sm')]: {
        maxHeight: '450px',
        maxWidth: '350px'
      }
    },
    dialogTitle: {
      fontWeight: 600,
      fontSize: '1.2rem'
    }
  };
});

export default SpecificContacts;
