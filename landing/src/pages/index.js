import React from "react";
import clsx from "clsx";
import Layout from "@theme/Layout";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import styles from "./index.module.css";
import HomepageFeatures from "../components/HomepageFeatures";

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={clsx("hero hero--primary", styles.heroBanner)}>
      <div className="container">
        <h1 className="hero__title">{siteConfig.title}</h1>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          {/* <Link
          
            className="button button--secondary button--lg"
            to="/docs/intro">
            Google Maps 📍
          </Link> */}
          <a
            target="_blank"
            href="https://goo.gl/maps/TvHWhvgrHLAi4vEQ8"
            className="button button--secondary button--lg"
          >
            Google Maps 📍
          </a>
        </div>
      </div>
    </header>
  );
}

export default function Home() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`Hello from ${siteConfig.title}`}
      description="Description will go into a meta tag in <head />"
    >
      <HomepageHeader />
      <main>
        <div
          style={{
            display: "flex",
            justifyContent: "space-around",
            marginTop: "2em",
            marginBottom: "2em",
          }}
        >
          <img src="/img/01-image.jpeg" style={{ width: "20vw" }} />
          <img src="/img/02-image.jpeg" style={{ width: "20vw" }} />
          <img src="/img/03-image.jpeg" style={{ width: "20vw" }} />
          <img src="/img/04-image.jpeg" style={{ width: "20vw" }} />
        </div>
        {/* <HomepageFeatures /> */}
      </main>
    </Layout>
  );
}
