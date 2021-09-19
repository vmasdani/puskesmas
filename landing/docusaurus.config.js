const lightCodeTheme = require("prism-react-renderer/themes/github");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  plugins: [
    [
      "docusaurus2-dotenv",
      {
        systemvars: true,
      },
    ],
  ],
  title: "Puskesmas Pasir Nangka",
  tagline: "Kecamatan Tigaraksa, Tangerang",
  url: "https://your-docusaurus-test-site.com",
  baseUrl: "/",
  // onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/puskesmas.png",
  organizationName: "facebook", // Usually your GitHub org/user name.
  projectName: "docusaurus", // Usually your repo name.
  themeConfig: {
    navbar: {
      title: "Puskesmas Pasir Nangka",
      logo: {
        alt: "My Site Logo",
        src: "img/puskesmas.png",
      },
      items: [
        // {
        //   type: "doc",
        //   docId: "intro",
        //   position: "left",
        //   label: "Tentang",
        // },
        // { to: "/blog", label: "Hubungi Kami", position: "left" },
        // { to: "/visimisi", label: "Visi & Misi", position: "left" },
        { to: "/gambaranumum", label: "Gambaran Umum", position: "left" },
        { to: "/pelayanan", label: "Pelayanan", position: "left" },

        
        {
          href: "/sidumas",
          // target: "_blank",
          label: "SIDUMAS",
          position: "left",
        },
        {
          href: "/admin",
          target: "_blank",
          label: "Berita",
          position: "left",
        },

        {
          href: "/admin",
          label: "Admin",
          position: "right",
        },

        // {
        //   href: "https://github.com/facebook/docusaurus",
        //   label: "GitHub",
        //   position: "right",
        // },
      ],
    },
    footer: {
      style: "dark",
      links: [
        {
          title: "Kontak",
          items: [
            {
              label: "WhatsApp +62 827-1829-3829",
              to: "#",
              target: "_blank",
            },
            {
              label: "Telp (021) 5001929320",
              to: "#",
              target: "_blank",
            },
          ],
        },
        {
          title: "Community",
          items: [
            {
              label: "Facebook",
              href: "https://stackoverflow.com/questions/tagged/docusaurus",
            },
            {
              label: "Instagram",
              href: "https://discordapp.com/invite/docusaurus",
            },
            // {
            //   label: "Twitter",
            //   href: "https://twitter.com/docusaurus",
            // },
          ],
        },
        {
          title: "Alamat",
          items: [
            {
              label:
                "Jl. Aria Jaya Santika No.Ds, Pasir Nangka, Tigaraksa Kec., Tangerang, Banten 15720",
              href: "#",
              // target: "_blank",
            },
            // {
            //   label: "GitHub",
            //   href: "https://github.com/facebook/docusaurus",
            // },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Puskesmas Pasir Nangka. Built with Docusaurus.`,
    },
    prism: {
      theme: lightCodeTheme,
      darkTheme: darkCodeTheme,
    },
  },
  presets: [
    [
      "@docusaurus/preset-classic",
      {
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          // Please change this to your repo.
          editUrl:
            "https://github.com/facebook/docusaurus/edit/master/website/",
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          editUrl:
            "https://github.com/facebook/docusaurus/edit/master/website/blog/",
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      },
    ],
  ],
};
