import type { SVGProps } from "react";

const LogoWithText = (props: SVGProps<SVGSVGElement>) => {
	return (
		<svg
			aria-label="Podevents — Podcast Notifications"
			fill="none"
			height="175"
			role="img"
			viewBox="26 44 708 175"
			width="auto"
			xmlns="http://www.w3.org/2000/svg"
			{...props}
		>
			<g>
				<path
					d="
      M78.5 69.5
      C88.7 53.8 105.2 45.8 121.5 45.8
      C143.3 45.8 162.2 58.7 171.1 78.4
      C177.2 92.0 176.7 105.9 171.8 120.0
      L148.8 184.2
      C141.7 204.3 126.7 214.3 106.5 215.0
      C121.1 207.4 132.8 193.7 138.7 176.9
      L161.4 112.9
      C165.7 100.8 165.9 91.3 161.8 82.0
      C155.1 66.8 140.8 57.0 124.0 56.7
      C111.5 56.5 99.2 61.6 90.8 71.4
      Z"
					fill="#70C667"
				/>

				<path
					d="
      M77.6 58.8
      C53.5 58.8 34.7 73.2 29.8 96.5
      C27.8 106.0 29.9 116.5 32.8 127.0
      L43.4 174.3
      C48.6 197.7 63.5 213.0 82.6 216.7
      C103.1 220.7 123.6 209.4 135.0 190.0
      C142.0 178.1 141.3 166.9 138.4 154.2
      L126.1 101.5
      C119.4 72.7 102.6 58.8 77.6 58.8
      Z"
					fill="#101622"
				/>

				<path
					d="M56.6 107.0L88.8 99.8"
					stroke="#70C667"
					strokeLinecap="round"
					strokeWidth="15.5"
				/>
				<path
					d="M63.1 138.7L95.4 131.5"
					stroke="#70C667"
					strokeLinecap="round"
					strokeWidth="15.5"
				/>
				<path
					d="M69.6 170.1L101.8 162.8"
					stroke="#70C667"
					strokeLinecap="round"
					strokeWidth="15.5"
				/>
			</g>

			<text
				fontFamily="Arial, Helvetica, sans-serif"
				fontSize="92"
				fontWeight="700"
				letterSpacing="-3"
				x="225"
				y="126"
			>
				<tspan fill="#FFFFFF">Pod</tspan>
				<tspan fill="#70C667">Events</tspan>
			</text>
			<text
				fill="#6E7A8B"
				fontFamily="Arial, Helvetica, sans-serif"
				fontSize="31"
				fontWeight="400"
				letterSpacing="5.2"
				x="229"
				y="174"
			>
				PODCAST NOTIFICATIONS
			</text>
		</svg>
	);
};

export default LogoWithText;
