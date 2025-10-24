import {
  initializeApp,
  type FirebaseApp,
  type FirebaseOptions,
} from "firebase/app";
import {
  createUserWithEmailAndPassword,
  getAuth,
  GithubAuthProvider,
  GoogleAuthProvider,
  onAuthStateChanged,
  signInWithEmailAndPassword,
  signInWithPopup,
  signOut,
  type Auth,
  type User,
} from "firebase/auth";
import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";

const firebaseConfig: FirebaseOptions = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
  projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
  storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
  appId: import.meta.env.VITE_FIREBASE_APP_ID,
  measurementId: import.meta.env.VITE_FIREBASE_MEASUREMENT_ID,
};

const authApp: FirebaseApp = initializeApp(firebaseConfig);
const Authorizer: Auth = getAuth(authApp);

const GoogleAuthorizer = new GoogleAuthProvider();

const GithubAuthorizer = new GithubAuthProvider();

interface FirebaseProviderProps {
  children: ReactNode;
}

const FirebaseContext = createContext<any | null>(null);

export const useFirebase = () => useContext(FirebaseContext);

export function FirebaseProvider(props: FirebaseProviderProps) {
  const [user, setUser] = useState<User | null>(null);

  const isLoggedIn = user ? true : false;

  useEffect(() => {
    onAuthStateChanged(Authorizer, (user) => {
      if (user) {
        setUser(user);
      } else {
        setUser(null);
      }
    });
  }, []);

  const SignUpWithEmailAndPassword = (email: string, password: string) =>
    createUserWithEmailAndPassword(Authorizer, email, password);

  const LogInWithEmailAndPassword = (email: string, password: string) =>
    signInWithEmailAndPassword(Authorizer, email, password);

  const LogInWithGoogleAccount = () => {
    signInWithPopup(Authorizer, GoogleAuthorizer);
  };

  const LogInWithGithubAccount = () => {
    signInWithPopup(Authorizer, GithubAuthorizer);
  };

  const LogOut = () => {
    signOut(Authorizer);
  };

  return (
    <FirebaseContext.Provider
      value={{
        SignUpWithEmailAndPassword,
        LogInWithEmailAndPassword,
        LogInWithGoogleAccount,
        LogInWithGithubAccount,
        isLoggedIn,
        LogOut,
      }}
    >
      {props.children}
    </FirebaseContext.Provider>
  );
}
