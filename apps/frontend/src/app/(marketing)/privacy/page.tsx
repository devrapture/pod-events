import type { Metadata } from "next";
import Link from "next/link";

export const metadata: Metadata = {
	title: "Privacy Policy",
	description:
		"How PodEvents collects, uses, and protects information when you use our podcast notification service.",
};

const LAST_UPDATED = "August 1, 2026";
const CONTACT_EMAIL = "devrapture@proton.me";

export default function PrivacyPage() {
	return (
		<section className="pt-28 pb-16 md:pt-32 md:pb-24">
			<article className="mx-auto max-w-3xl px-6">
				<header className="mb-12">
					<p className="mb-3 text-sm text-zinc-500">Legal</p>
					<h1 className="font-semibold text-3xl text-zinc-50 tracking-tight md:text-4xl">
						Privacy Policy
					</h1>
					<p className="mt-4 text-zinc-400 leading-relaxed">
						This page describes how PodEvents handles information when you use
						our website and podcast notification service. It is written in plain
						language for transparency and is not legal advice.
					</p>
					<p className="mt-3 text-sm text-zinc-500">
						Last updated: {LAST_UPDATED}
					</p>
				</header>

				<div className="space-y-10 text-zinc-400 leading-relaxed">
					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							1. Information we collect
						</h2>
						<p className="mb-3">
							Depending on how you use PodEvents, we may collect the following
							categories of information:
						</p>
						<ul className="list-disc space-y-2 pl-5">
							<li>
								<span className="text-zinc-300">Account and profile data</span>{" "}
								from Spotify when you sign in, such as your name, email address,
								profile avatar URL, and Spotify user ID.
							</li>
							<li>
								<span className="text-zinc-300">
									Authentication credentials
								</span>{" "}
								needed to keep your session working, including Spotify access
								and refresh tokens (stored encrypted) and a PodEvents session
								token.
							</li>
							<li>
								<span className="text-zinc-300">Podcast preferences</span> you
								create in the product, such as shows you subscribe to monitor
								and related show metadata (name, description, image, Spotify
								URL, and latest episode tracking information).
							</li>
							<li>
								<span className="text-zinc-300">Notification settings</span> you
								configure, such as Slack or Discord webhook URLs, Telegram chat
								IDs, channel labels, and whether a channel is active.
							</li>
							<li>
								<span className="text-zinc-300">
									Technical and operational data
								</span>{" "}
								generated while providing the service, such as authentication
								state, API request activity needed to operate the product, and
								notification delivery logs used to operate and troubleshoot
								alerts.
							</li>
						</ul>
						<p className="mt-3">
							We do not require you to provide information beyond what is needed
							to authenticate, connect services, and deliver the notifications
							you configure.
						</p>
					</section>

					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							2. How we use information
						</h2>
						<p className="mb-3">We use the information described above to:</p>
						<ul className="list-disc space-y-2 pl-5">
							<li>Create and maintain your PodEvents account</li>
							<li>Authenticate you and keep your session secure</li>
							<li>
								Access Spotify on your behalf (within the permissions you grant)
								to find shows and monitor new episodes
							</li>
							<li>
								Send new-episode notifications to the channels you connect
							</li>
							<li>
								Operate, secure, debug, and improve the reliability of the
								service
							</li>
							<li>
								Respond to support requests and communicate about the service
								when needed
							</li>
						</ul>
						<p className="mt-3">We do not sell your personal information.</p>
					</section>

					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							3. Cookies and similar technologies
						</h2>
						<p className="mb-3">
							PodEvents uses a limited set of cookies and browser storage
							technologies that are necessary for authentication and routing:
						</p>
						<ul className="list-disc space-y-2 pl-5">
							<li>
								A session cookie (
								<code className="rounded bg-white/[0.06] px-1.5 py-0.5 text-sm text-zinc-300">
									pod_events_token
								</code>
								) that mirrors your signed-in state so protected routes can work
								correctly
							</li>
							<li>
								Browser{" "}
								<code className="rounded bg-white/[0.06] px-1.5 py-0.5 text-sm text-zinc-300">
									localStorage
								</code>{" "}
								for storing your PodEvents session token
							</li>
							<li>
								Short-lived OAuth state values used during Spotify sign-in to
								help prevent CSRF attacks
							</li>
						</ul>
						<p className="mt-3">
							PodEvents uses{" "}
							<a
								className="text-emerald-400 transition-colors hover:text-emerald-300"
								href="https://www.sabilytics.com/"
								rel="noopener noreferrer"
								target="_blank"
							>
								Sabilytics
							</a>{" "}
							for website usage analytics. Sabilytics may process technical and
							usage information when you visit the site. We do not use it for
							third-party advertising.
						</p>
					</section>

					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							4. Authentication and connected third-party services
						</h2>
						<p className="mb-3">
							PodEvents relies on third-party services you choose to connect:
						</p>
						<ul className="list-disc space-y-2 pl-5">
							<li>
								<span className="text-zinc-300">Spotify</span> — used for
								sign-in and podcast library access via OAuth. Spotify may
								process your data under its own privacy policy.
							</li>
							<li>
								<span className="text-zinc-300">Slack and Discord</span> — when
								you add webhook destinations, we send episode notifications to
								those endpoints.
							</li>
							<li>
								<span className="text-zinc-300">Telegram</span> — when you
								connect a Telegram chat, we use Telegram APIs/webhooks to
								deliver notifications.
							</li>
						</ul>
						<p className="mt-3">
							Those providers process information under their own terms and
							privacy policies. PodEvents only uses the connections you
							configure to provide the features you request.
						</p>
					</section>

					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							5. Music, podcast, and listening-related data
						</h2>
						<p className="mb-3">
							PodEvents is focused on podcast show monitoring, not full music
							playback history. With your Spotify authorization, we may access:
						</p>
						<ul className="list-disc space-y-2 pl-5">
							<li>
								Profile information needed for account creation (including email
								and basic profile fields)
							</li>
							<li>
								Saved shows and followed podcasts so you can import or subscribe
								to shows you care about
							</li>
							<li>
								Show and episode metadata needed to detect new releases and
								compose notifications
							</li>
						</ul>
						<p className="mt-3">
							We store the podcast subscriptions and notification preferences
							you set up in PodEvents. We do not use your Spotify connection to
							build advertising profiles.
						</p>
					</section>

					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							6. Data sharing and service providers
						</h2>
						<p className="mb-3">
							We share information only as needed to operate PodEvents:
						</p>
						<ul className="list-disc space-y-2 pl-5">
							<li>
								With connected notification platforms you configure (for
								example, posting a new-episode alert to a Slack webhook)
							</li>
							<li>
								With Spotify as required to authenticate and retrieve podcast
								data under the scopes you grant
							</li>
							<li>
								With infrastructure providers that host or process the
								application and database when you use a hosted deployment of
								PodEvents
							</li>
							<li>
								If required by law, legal process, or to protect the security
								and integrity of the service
							</li>
						</ul>
						<p className="mt-3">
							Because PodEvents is open source, a self-hosted deployment is
							controlled by the operator of that instance. If you use a
							self-hosted copy, that operator&apos;s practices apply to the data
							processed there.
						</p>
					</section>

					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							7. Data retention
						</h2>
						<p>
							We retain account, subscription, channel, token, and operational
							records for as long as needed to provide the service, maintain
							security, and meet legitimate operational needs. Exact retention
							periods are not fixed in the product code. If you delete or stop
							using connected channels, or if account-deletion tooling becomes
							available, related data will be removed or disabled according to
							that process. You can also disconnect third-party services
							(including Spotify authorizations) through those providers&apos;
							settings.
						</p>
					</section>

					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							8. Data security
						</h2>
						<p className="mb-3">
							We design PodEvents with common production security practices,
							including:
						</p>
						<ul className="list-disc space-y-2 pl-5">
							<li>OAuth 2.0 for Spotify authentication</li>
							<li>AES-256-GCM encryption for stored Spotify credentials</li>
							<li>JWT-based session authentication for API access</li>
							<li>
								Secure cookie attributes for the session cookie where supported
								(for example,{" "}
								<code className="rounded bg-white/[0.06] px-1.5 py-0.5 text-sm text-zinc-300">
									SameSite=Lax
								</code>{" "}
								and{" "}
								<code className="rounded bg-white/[0.06] px-1.5 py-0.5 text-sm text-zinc-300">
									Secure
								</code>
								)
							</li>
						</ul>
						<p className="mt-3">
							No method of transmission or storage is completely secure. We work
							to protect your information, but we cannot guarantee absolute
							security.
						</p>
					</section>

					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							9. Your rights and choices
						</h2>
						<p className="mb-3">
							Depending on where you live, you may have rights to access,
							correct, delete, or restrict certain processing of your personal
							information, or to object to processing and request portability.
							You can also:
						</p>
						<ul className="list-disc space-y-2 pl-5">
							<li>Sign out of PodEvents to clear your local session token</li>
							<li>
								Remove notification channels you no longer want PodEvents to use
							</li>
							<li>
								Manage or revoke Spotify app access from your Spotify account
								settings
							</li>
							<li>
								Contact us to request help with access, correction, or deletion
								of account data associated with the service we operate
							</li>
						</ul>
						<p className="mt-3">
							We will respond to requests in accordance with applicable law. We
							may need to verify your identity before fulfilling certain
							requests.
						</p>
					</section>

					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							10. Children&apos;s privacy
						</h2>
						<p>
							PodEvents is not directed to children under 13 (or the equivalent
							minimum age in your jurisdiction), and we do not knowingly collect
							personal information from children. If you believe a child has
							provided personal information to us, please contact us so we can
							take appropriate steps.
						</p>
					</section>

					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							11. International data processing
						</h2>
						<p>
							PodEvents may be accessed from different countries. If you use a
							hosted deployment of the service, your information may be
							processed in locations other than where you live. Those locations
							may have data-protection laws that differ from the laws in your
							jurisdiction. By using the service, you understand that your
							information may be transferred and processed in this way as needed
							to provide PodEvents.
						</p>
					</section>

					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							12. Changes to this policy
						</h2>
						<p>
							We may update this Privacy Policy from time to time. When we do,
							we will change the &quot;Last updated&quot; date on this page. If
							changes are material, we may also provide additional notice
							through the product or site. Continued use of PodEvents after an
							update means you acknowledge the revised policy.
						</p>
					</section>

					<section>
						<h2 className="mb-3 font-semibold text-xl text-zinc-50 tracking-tight">
							13. Contact us
						</h2>
						<p className="mb-3">
							If you have questions about this Privacy Policy or about how
							PodEvents handles personal information, contact us at:
						</p>
						<p>
							<a
								className="text-emerald-400 transition-colors hover:text-emerald-300"
								href={`mailto:${CONTACT_EMAIL}`}
							>
								{CONTACT_EMAIL}
							</a>
						</p>
						<p className="mt-6">
							<Link
								className="text-sm text-zinc-500 transition-colors hover:text-zinc-300"
								href="/"
							>
								← Back to home
							</Link>
						</p>
					</section>
				</div>
			</article>
		</section>
	);
}
