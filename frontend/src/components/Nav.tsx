import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

type NavProperties = {
  jwt: string;
  setJwt: React.Dispatch<React.SetStateAction<string>>;
};

const Nav = ({ jwt, setJwt }: NavProperties) => {
  const [loginLink, setLoginLink] = useState(<></>);

  const logout = () => {
    setJwt('');
    window.localStorage.removeItem('jwt');
  };

  useEffect(() => {
    if (jwt == '') {
      setLoginLink(
        <ul className='navbar-nav me-auto mb-2 mb-lg-0'>
          <li className='nav-item'>
            <Link to='/signup' className='nav-link'>
              Signup
            </Link>
          </li>
          <li className='nav-item'>
            <Link to='/login' className='nav-link'>
              Login
            </Link>
          </li>
        </ul>,
      );
    } else {
      setLoginLink(
        <Link to='/logout' className='nav-link' onClick={logout}>
          Logout
        </Link>,
      );
    }
  }, [jwt]);

  return (
    <nav className='navbar navbar-expand-lg bg-secondary navbar-dark bg-opacity-75'>
      <div className='container-fluid'>
        <Link to='/' className='navbar-brand'>
          Touring log
        </Link>
        <button
          className='navbar-toggler'
          type='button'
          data-bs-toggle='collapse'
          data-bs-target='#navbarTogglerDemo02'
          aria-controls='navbarTogglerDemo02'
          aria-expanded='false'
          aria-label='Toggle navigation'
        >
          <span className='navbar-toggler-icon'></span>
        </button>
        <div className='collapse navbar-collapse' id='navbarSupportedContent'>
          <ul className='navbar-nav me-auto mb-2 mb-lg-0'>
            <li className='nav-item'>
              <Link to='/' className='nav-link'>
                Home
              </Link>
            </li>
            <li className='nav-item'>
              <Link to='/upload' className='nav-link'>
                Upload
              </Link>
            </li>
          </ul>
          <div className='text-end'>{loginLink}</div>
        </div>
      </div>
    </nav>
  );
};
export default Nav;
