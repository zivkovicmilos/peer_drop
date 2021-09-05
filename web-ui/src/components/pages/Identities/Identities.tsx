import { Box } from '@material-ui/core';
import AddRoundedIcon from '@material-ui/icons/AddRounded';
import { FC, useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import Pagination from '../../atoms/Pagination/Pagination';
import usePagination from '../../atoms/Pagination/pagination.hook';
import IdentityCard from '../../molecules/IdentityCard/IdentityCard';
import IdentitySort from '../../molecules/IdentitySort/IdentitySort';
import {
  EIdentitySortDirection,
  EIdentitySortParam
} from '../../molecules/IdentitySort/identitySort.types';
import { IIdentitiesProps, IIdentity } from './identities.types';

const Identities: FC<IIdentitiesProps> = () => {
  const history = useHistory();

  const [activeSort, setActiveSort] = useState<EIdentitySortParam>(
    EIdentitySortParam.NAME
  );
  const [sortDirection, setSortDirection] = useState<EIdentitySortDirection>(
    EIdentitySortDirection.ASC
  );

  const handleNewContact = () => {
    history.push('/identities/new');
  };

  const { page, count, setCount, limit, handlePageChange } = usePagination({
    limit: 8
  });

  // TODO modify this page change listener
  useEffect(() => {
    setCount(identities.length);
  }, []);

  // TODO add useEffect listeners for the sort params / direction

  // TODO add useEffect for fetching the list of identities

  const identities: IIdentity[] = [
    {
      id: '1',
      picture: '1',
      name: 'Milos',

      publicKeyID: '4AEE18F83AFDEB23',
      numWorkspaces: 4,
      creationDate: '20.08.2021.'
    },
    {
      id: '2',
      picture: '1',
      name: 'Milos',

      publicKeyID: '4AEE18F83AFDEB23',
      numWorkspaces: 4,
      creationDate: '20.08.2021.'
    },
    {
      id: '2',
      picture: '1',
      name: 'Milos',

      publicKeyID: '4AEE18F83AFDEB23',
      numWorkspaces: 4,
      creationDate: '20.08.2021.'
    },
    {
      id: '3',
      picture: '1',
      name: 'Milos',

      publicKeyID: '4AEE18F83AFDEB23',
      numWorkspaces: 4,
      creationDate: '20.08.2021.'
    },
    {
      id: '4',
      picture: '1',
      name: 'Milos',

      publicKeyID: '4AEE18F83AFDEB23',
      numWorkspaces: 4,
      creationDate: '20.08.2021.'
    }
  ];

  return (
    <Box
      display={'flex'}
      flexDirection={'column'}
      width={'100%'}
      height={'100%'}
    >
      <Box
        display={'flex'}
        justifyContent={'space-between'}
        width={'100%'}
        alignItems={'center'}
        mb={4}
      >
        <PageTitle title={'Identities'} />
        <ActionButton
          text={'New identity'}
          startIcon={<AddRoundedIcon />}
          onClick={handleNewContact}
        />
      </Box>

      <Box display={'flex'} width={'100%'} mb={4}>
        <IdentitySort
          activeSort={activeSort}
          setActiveSort={setActiveSort}
          setSortDirection={setSortDirection}
          sortDirection={sortDirection}
        />
      </Box>

      <Box display={'flex'} width={'100%'} flexDirection={'column'}>
        <Box
          display={'flex'}
          flexWrap={'wrap'}
          width={'100%'}
          ml={'-36px'}
          mt={'-36px'}
        >
          {identities.map((identity) => {
            return (
              <IdentityCard
                picture={identity.picture}
                name={identity.name}
                publicKeyID={identity.publicKeyID}
                numWorkspaces={identity.numWorkspaces}
                creationDate={identity.creationDate}
              />
            );
          })}
        </Box>
        <Pagination
          count={count}
          limit={limit}
          page={page}
          onPageChange={handlePageChange}
        />
      </Box>
    </Box>
  );
};

export default Identities;
