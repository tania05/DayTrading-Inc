module Main where

import Control.Applicative ((<$>))
import Control.Concurrent (forkIO, threadDelay)
import Control.Exception (bracket_)
import Control.Monad (replicateM)
import Data.Char as C
import Data.Maybe
import Data.Time.Clock (getCurrentTime)
import Data.Time.Clock.POSIX
import Data.Time.Format (defaultTimeLocale, formatTime, iso8601DateFormat)
import Data.List (intercalate)
import GHC.IO.Handle
import Network
import System.IO.Error (catchIOError)
import System.IO (hPutStrLn)
import System.Random (getStdRandom, randomR, randomRIO)
import Text.Printf (printf)

data Command = QuoteCommand
  { stockSym :: String
  , userid :: String
  } deriving (Show)

main :: IO ()
main = listenOn port >>= listenLoop
  where
    port = PortNumber 8080

strip :: String -> String
strip str = reverse $ dropWhile C.isSpace (reverse $ dropWhile C.isSpace str)

listenLoop :: Socket -> IO ()
listenLoop sock = do
  (handler, _, _) <- accept sock
  forkIO $ handleConnection handler
  listenLoop sock
  where
    getHandle (handler, _, _) = handler

handleConnection :: Handle -> IO ()
handleConnection conn =
  bracket_
    (return ())
    (do hClose conn
        putStrLn "Connection closed")
    (handleConnection' conn)

handleConnection' :: Handle -> IO ()
handleConnection' conn = do
  command <- hGetLine conn >>= \c -> return $ parseCommand c
  giveResponse conn command

parseCommand :: String -> Maybe Command
parseCommand [] = Nothing
parseCommand command =
  parsePart (takeWhile (/= ',') command) (tail $ dropWhile (/= ',') command)
  where
    parsePart :: String -> String -> Maybe Command
    parsePart rawSym rawUserId
      | length sym > 3 = Nothing
      | any (not . C.isPrint) sym = Nothing
      | any (not . C.isPrint) userId = Nothing
      | otherwise = Just $ QuoteCommand sym userId
      where
        userId = strip rawUserId
        sym = strip rawSym

giveResponse :: Handle -> Maybe Command -> IO ()
giveResponse _ Nothing = putStrLn "Failed to parse input"
giveResponse handle (Just (QuoteCommand sym userid)) = do
  quote <- quoteIO
  timestamp <- timestampIO
  cryptokey <- cryptokeyIO
  let output = intercalate "," [quote, sym, userid, timestamp, cryptokey]
  delayTime <- randomRIO (2500000, 10000000) :: IO Int
  threadDelay delayTime
  hPutStrLn handle output
  putStrLn $ "Quote lookup: " ++ show (QuoteCommand sym userid) ++ " returned " ++ output
  where
    quoteIO :: IO String
    quoteIO = do
      num <- randomRIO (1, 100000) :: IO Integer
      return $ printf "%d.%02d" (num `quot` 100) (num `mod` 100)
    timestampIO :: IO String
    timestampIO = do
      time <- getPOSIXTime
      return $ show $ round (time * 1000)
    cryptokeyIO :: IO String
    cryptokeyIO =
      replicateM 25 $ fmap (\x -> chr $ x + 97) (getStdRandom $ randomR (0, 25))
