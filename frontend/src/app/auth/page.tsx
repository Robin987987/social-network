"use client"

import Head from 'next/head'
import Header from '../../components/headers/LoginHeader'
import RegisterForm from '../../components/auth/Register'
import background from '../../../public/assets/background.webp';


export default function Auth()  {
    return (

        <div className="flex flex-col h-screen">
            <Header />
            <Head>
                <title>IrieSphere</title>
            </Head>
            <div><RegisterForm />
            </div>
            
                
            </div>

    )
}