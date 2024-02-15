#!/bin/bash

MOCK_MAIL_SENDER=no-reply@mail.flyover.rifcomputing.net

awslocal ses verify-email-identity --email $MOCK_MAIL_SENDER