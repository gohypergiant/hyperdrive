import logging
import os
import signal
import sys
import time
import traceback
from dataclasses import dataclass
from logging import handlers

from .__doc__ import WalkerDoc
from .signal_handlers import WalkerSignalHandlers
from .methods import WalkerMethods


class DaemonError(Exception):
    """Base exception class for errors from this module."""


class DaemonProcessDetachError(DaemonError, OSError):
    """Exception raised when process detach fails."""


@dataclass
class Walker(WalkerSignalHandlers, WalkerMethods, WalkerDoc):
    walker_interval: int = 15
    service_dir: str = "/home/jovyan/.executor"
    job_dir: str = "/home/jovyan/_jobs"
    log_level = logging.DEBUG
    log_name: str = "events.log"
    log_prune: int = 5  # days

    def __post_init__(self):

        os.makedirs(self.service_dir, exist_ok=True)
        self.logger = self._initialize_logger()

        self.logger.info(
            f"Initializing Walker Service with service interval: {self.walker_interval}"
        )
        full_log_path = self.service_dir + "/" + self.log_name
        self.logger.info(f"\tLogs written to: {full_log_path}")

        signal.signal(signal.SIGQUIT, self.signal_quit_handler)
        signal.signal(signal.SIGTERM, self.signal_term_handler)
        signal.signal(signal.SIGHUP, self.signal_hup_handler)

        self.main()

    def fork_then_exit_parent(self, error_message):
        """Fork a child process, then exit the parent process.

        :param error_message: Message for the exception in case of a
            detach failure.
        :return: ``None``.
        :raise DaemonProcessDetachError: If the fork fails.
        """
        try:
            pid = os.fork()
            if pid > 0:
                os._exit(0)
        except OSError as exc:
            error = DaemonProcessDetachError(
                "{message}: [{exc.errno:d}] {exc.strerror}".format(
                    message=error_message, exc=exc
                )
            )
            raise error from exc

    def main(self):
        self.fork_then_exit_parent(error_message="Failed first fork")
        os.setsid()
        self.fork_then_exit_parent(error_message="Failed second fork")
        walk = True
        try:
            while walk:
                self.logger.info("Running scheduler.")
                next_job_name = self.walk_job_tree()
                if next_job_name is None:
                    self.logger.info("END: No more jobs to run.")
                else:
                    self.logger.info(f"Executing job: {next_job_name}")
                    self.run_next_job(next_job_name)
                    time.sleep(self.walker_interval)
                time.sleep(self.walker_interval)
        except:
            self.logger.critical(traceback.format_exc())
            sys.exit(1)

        sys.exit(0)

    def _initialize_logger(self):
        formatter = logging.Formatter(
            "%(asctime)s.%(msecs)03d > %(levelname)s > %(message)s", "%Y-%m-%d %H:%M:%S"
        )
        logger = logging.getLogger(self.log_name)
        logger.setLevel(self.log_level)
        logHandler = handlers.TimedRotatingFileHandler(
            os.path.join(self.service_dir, self.log_name),
            when="D",
            backupCount=self.log_prune,
        )
        logHandler.setFormatter(formatter)
        logger.addHandler(logHandler)
        return logger


__doc__ = WalkerDoc.__doc__
