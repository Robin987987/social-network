"use client"

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
// Your authentication hook or context
import { useAuth } from '../util/utils';
import MainHeader from '@/components/headers/MainHeader';
import LeftNavBar from '@/components/leftnavbar/LeftNavBar'
import background from '../../public/assets/background.png';
import PostFeed from '@/components/postfeed/PostFeed';

export default function Home() {
    const router = useRouter();
    useEffect(() => {
        (async () => {
            const auth = await useAuth();
            if (!auth.is_authenticated) {
                router.push('/auth');
            }
        })();
    }, []); // Empty dependency array to run only once on mount

    return (

        <div>
            {/* Header */}
            <header>
                <MainHeader />
            </header>

            {/* Main Content */}
            <main>
                <section>
                    <PostFeed />
                </section>

                {/* Sidebar */}
                <aside>
                    {/* Display friend suggestions */}
                </aside>
                
            </main>

            {/* Footer */}
            <footer>
                {/* Display copyright information */}
            </footer>
            
        </div>
    );
}
