import { Box } from '@material-ui/core';
import { FC } from 'react';
import SuggestedFile from '../../atoms/SuggestedFile/SuggestedFile';
import { ISuggestedListProps } from './suggestedList.types';

const SuggestedList: FC<ISuggestedListProps> = (props) => {
  const { files } = props;

  return (
    <Box display={'flex'} width={'80%'} justifyContent={'space-between'}>
      {files.map((file) => {
        return <SuggestedFile file={file} />;
      })}
    </Box>
  );
};

export default SuggestedList;