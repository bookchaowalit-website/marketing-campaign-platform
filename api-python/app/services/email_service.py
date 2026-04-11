"""
Email service for Python API (Marketing emails, notifications)
Supports SendGrid and AWS SES
"""
import os
import logging
from typing import Optional
from dataclasses import dataclass

logger = logging.getLogger(__name__)


@dataclass
class EmailConfig:
    provider: str = "stub"  # stub, sendgrid, ses
    api_key: Optional[str] = None  # SendGrid API key or AWS access key
    api_secret: Optional[str] = None  # AWS secret key
    region: str = "us-east-1"
    from_email: str = "noreply@marketing-saas.com"
    from_name: str = "Marketing SaaS"


class EmailService:
    """Email service for sending marketing emails"""

    def __init__(self, config: EmailConfig = None):
        self.config = config or EmailConfig(
            provider=os.getenv("EMAIL_PROVIDER", "stub"),
            api_key=os.getenv("SENDGRID_API_KEY") or os.getenv("AWS_ACCESS_KEY_ID"),
            api_secret=os.getenv("AWS_SECRET_ACCESS_KEY"),
            region=os.getenv("AWS_REGION", "us-east-1"),
            from_email=os.getenv("EMAIL_FROM_ADDRESS", "noreply@marketing-saas.com"),
            from_name=os.getenv("EMAIL_FROM_NAME", "Marketing SaaS"),
        )
        self._client = None

    async def send_email(self, to: str, subject: str, html_content: str) -> bool:
        """Send an email"""
        if self.config.provider == "sendgrid":
            return await self._send_sendgrid(to, subject, html_content)
        elif self.config.provider == "ses":
            return await self._send_ses(to, subject, html_content)
        else:
            return self._send_stub(to, subject, html_content)

    async def _send_sendgrid(self, to: str, subject: str, html_content: str) -> bool:
        """Send email via SendGrid"""
        try:
            import sendgrid
            from sendgrid.helpers.mail import Mail

            sg = sendgrid.SendGridAPIClient(api_key=self.config.api_key)
            message = Mail(
                from_email=(self.config.from_email, self.config.from_name),
                to_emails=to,
                subject=subject,
                html_content=html_content,
            )
            response = sg.send(message)
            logger.info(f"Email sent to {to} via SendGrid: {response.status_code}")
            return 200 <= response.status_code < 300
        except Exception as e:
            logger.error(f"Failed to send email via SendGrid: {e}")
            return False

    async def _send_ses(self, to: str, subject: str, html_content: str) -> bool:
        """Send email via AWS SES"""
        try:
            import boto3

            ses = boto3.client(
                "ses",
                aws_access_key_id=self.config.api_key,
                aws_secret_access_key=self.config.api_secret,
                region_name=self.config.region,
            )
            ses.send_email(
                Source=f"{self.config.from_name} <{self.config.from_email}>",
                Destination={"ToAddresses": [to]},
                Message={
                    "Subject": {"Data": subject, "Charset": "UTF-8"},
                    "Body": {
                        "Html": {"Data": html_content, "Charset": "UTF-8"},
                    },
                },
            )
            logger.info(f"Email sent to {to} via SES")
            return True
        except Exception as e:
            logger.error(f"Failed to send email via SES: {e}")
            return False

    def _send_stub(self, to: str, subject: str, html_content: str) -> bool:
        """Stub implementation for development"""
        logger.info(f"[STUB] Email to {to}: {subject}")
        # Extract URL from HTML for console output
        import re
        urls = re.findall(r'href=[\'"]?([^\'" >]+)', html_content)
        if urls:
            for url in urls[:2]:  # Print first 2 URLs
                logger.info(f"[STUB] URL: {url}")
        return True

    async def send_password_reset(self, to: str, reset_url: str) -> bool:
        """Send password reset email"""
        html = self._password_reset_template(reset_url)
        return await self.send_email(to, "Reset Your Password", html)

    async def send_verification(self, to: str, verification_url: str) -> bool:
        """Send email verification email"""
        html = self._verification_template(verification_url)
        return await self.send_email(to, "Verify Your Email", html)

    async def send_marketing_newsletter(self, to: str, content: dict) -> bool:
        """Send marketing newsletter email"""
        html = self._newsletter_template(content)
        return await self.send_email(to, content.get("subject", "Marketing Update"), html)

    def _password_reset_template(self, reset_url: str) -> str:
        return f"""
        <!DOCTYPE html>
        <html>
        <head>
            <meta charset="UTF-8">
            <style>
                body {{ font-family: Arial, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }}
                .container {{ max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; overflow: hidden; }}
                .header {{ background: linear-gradient(135deg, #667eea, #764ba2); color: white; padding: 40px 20px; text-align: center; }}
                .content {{ padding: 40px 30px; }}
                .button {{ display: inline-block; background: #667eea; color: white; text-decoration: none; padding: 14px 32px; border-radius: 6px; margin: 20px 0; }}
                .footer {{ background: #f8f9fa; padding: 20px; text-align: center; color: #666; font-size: 12px; }}
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header"><h1>🔑 Reset Your Password</h1></div>
                <div class="content">
                    <p>Click the button below to reset your password:</p>
                    <p style="text-align: center;"><a href="{reset_url}" class="button">Reset Password</a></p>
                    <p>This link expires in <strong>1 hour</strong>.</p>
                    <p style="font-size: 12px; word-break: break-all;">{reset_url}</p>
                </div>
                <div class="footer"><p>&copy; 2024 Marketing SaaS</p></div>
            </div>
        </body>
        </html>
        """

    def _verification_template(self, verification_url: str) -> str:
        return f"""
        <!DOCTYPE html>
        <html>
        <head>
            <meta charset="UTF-8">
            <style>
                body {{ font-family: Arial, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }}
                .container {{ max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; overflow: hidden; }}
                .header {{ background: linear-gradient(135deg, #11998e, #38ef7d); color: white; padding: 40px 20px; text-align: center; }}
                .content {{ padding: 40px 30px; }}
                .button {{ display: inline-block; background: #11998e; color: white; text-decoration: none; padding: 14px 32px; border-radius: 6px; margin: 20px 0; }}
                .footer {{ background: #f8f9fa; padding: 20px; text-align: center; color: #666; font-size: 12px; }}
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header"><h1>✅ Verify Your Email</h1></div>
                <div class="content">
                    <p>Welcome! Please verify your email:</p>
                    <p style="text-align: center;"><a href="{verification_url}" class="button">Verify Email</a></p>
                    <p>This link expires in <strong>24 hours</strong>.</p>
                    <p style="font-size: 12px; word-break: break-all;">{verification_url}</p>
                </div>
                <div class="footer"><p>&copy; 2024 Marketing SaaS</p></div>
            </div>
        </body>
        </html>
        """

    def _newsletter_template(self, content: dict) -> str:
        title = content.get("title", "Marketing Update")
        body = content.get("body", "")
        cta_url = content.get("cta_url", "")
        cta_text = content.get("cta_text", "Learn More")

        return f"""
        <!DOCTYPE html>
        <html>
        <head>
            <meta charset="UTF-8">
            <style>
                body {{ font-family: Arial, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }}
                .container {{ max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; overflow: hidden; }}
                .header {{ background: linear-gradient(135deg, #f093fb, #f5576c); color: white; padding: 40px 20px; text-align: center; }}
                .content {{ padding: 40px 30px; }}
                .button {{ display: inline-block; background: #f5576c; color: white; text-decoration: none; padding: 14px 32px; border-radius: 6px; margin: 20px 0; }}
                .footer {{ background: #f8f9fa; padding: 20px; text-align: center; color: #666; font-size: 12px; }}
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header"><h1>📧 {title}</h1></div>
                <div class="content">
                    <p>{body}</p>
                    {f'<p style="text-align: center;"><a href="{cta_url}" class="button">{cta_text}</a></p>' if cta_url else ''}
                </div>
                <div class="footer"><p>&copy; 2024 Marketing SaaS | <a href="{{unsubscribe_url}}">Unsubscribe</a></p></div>
            </div>
        </body>
        </html>
        """


# Singleton instance
_email_service: Optional[EmailService] = None


def get_email_service() -> EmailService:
    """Get or create email service singleton"""
    global _email_service
    if _email_service is None:
        _email_service = EmailService()
    return _email_service
