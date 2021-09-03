import { Box } from '@material-ui/core';
import { FC } from 'react';
import { useParams } from 'react-router-dom';
import { IContactEditParams, IContactEditProps } from './contactEdit.types';

const ContactEdit: FC<IContactEditProps> = (props) => {
  const { type } = props;

  const { contactId } = useParams() as IContactEditParams;

  return (
    <Box>
      Contact {type}! {contactId}
    </Box>
  );
};

export default ContactEdit;
