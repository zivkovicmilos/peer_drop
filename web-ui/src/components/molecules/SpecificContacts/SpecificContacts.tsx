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
import ContactsService from '../../../services/contacts/contactsService';
import { IContactResponse } from '../../../services/contacts/contactsService.types';
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

  const [allContacts, setAllContacts] = useState<{ data: IContactResponse[] } | undefined>();

  const [filteredList, setFilteredList] = useState<{ data: IContactResponse[] } | undefined>({ data: [] });

  const [selectedContacts, setSelectedContacts] = useState<{
    data: IContactResponse[];
  }>({
    data: contactIDs.contacts
  });

  useEffect(() => {
    if (dialogOpen) {
      const fetchAllContacts = async () => {
        return await ContactsService.getContacts({ page: -1, limit: -1 });
      };

      fetchAllContacts()
        .then((response) => {
          setAllContacts({ data: response.data });

          setFilteredList({ data: response.data });
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
    return selectedContacts.data.some((contact) => contact.id === id);
  };

  const findContact = (id: string) => {
    for (let i = 0; i < selectedContacts.data.length; i++) {
      if (selectedContacts.data[i].id === id) return i;
    }

    return -1;
  };

  const handleContactToggle = async (data: IContactResponse) => {
    const index = findContact(data.id);

    if (index > -1) {
      let newArr = selectedContacts.data;
      newArr.splice(index, 1);

      setSelectedContacts({ data: newArr });
    } else {
      let newArr = selectedContacts.data;
      newArr.push(data);

      setSelectedContacts({ data: newArr });
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

    let tempFilteredContacts = { data: allContacts.data };
    // Filter if there is a value in search input
    if (value) {
      tempFilteredContacts.data = manualFilter(
        (contact: any) =>
          contact.name.toLowerCase().includes(value.toLowerCase()) ||
          contact.publicKeyID.toLowerCase().includes(value.toLowerCase()),
        allContacts.data
      );
    }

    setFilteredList(tempFilteredContacts);
  };

  const handleDialogToggle = (open: boolean) => {
    setDialogOpen(open);

    let contactIDs: IContactResponse[] = [];
    for (let i = 0; i < selectedContacts.data.length; i++) {
      contactIDs.push(selectedContacts.data[i]);
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
        {selectedContacts.data.length < 1 ? (
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
            {selectedContacts.data.map(
              (selectedContact: IContactResponse, index) => {
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
          filteredList.data.map((contact) => {
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
                onChange={() => handleContactToggle(contact)}
              />
            );
          })}
          {loading && <LoadingIndicator size={54} />}
          {!loading && filteredList && filteredList.data.length < 1 && (
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
